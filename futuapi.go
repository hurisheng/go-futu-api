package futuapi

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/hurisheng/go-futu-api/pb/common"
	"github.com/hurisheng/go-futu-api/pb/getdelaystatistics"
	"github.com/hurisheng/go-futu-api/pb/getglobalstate"
	"github.com/hurisheng/go-futu-api/pb/getuserinfo"
	"github.com/hurisheng/go-futu-api/pb/initconnect"
	"github.com/hurisheng/go-futu-api/pb/keepalive"
	"github.com/hurisheng/go-futu-api/pb/notify"
	"github.com/hurisheng/go-futu-api/pb/verification"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDInitConnect        = 1001 //InitConnect	初始化连接
	ProtoIDGetGlobalState     = 1002 //GetGlobalState	获取全局状态
	ProtoIDNotify             = 1003 //Notify	系统通知推送
	ProtoIDKeepAlive          = 1004 //KeepAlive	保活心跳
	ProtoIDGetUserInfo        = 1005 //GetUserInfo 获取用户信息
	ProtoIDVerification       = 1006 // 请求或输入验证码
	ProtoIDGetDelayStatistics = 1007 //GetDelayStatistics  获取延迟统计
)

var workers map[uint32]protocol.Worker = make(map[uint32]protocol.Worker)

func init() {
	workers[ProtoIDInitConnect] = protocol.NewGetter()
	workers[ProtoIDGetGlobalState] = protocol.NewGetter()
	workers[ProtoIDNotify] = protocol.NewUpdater()
	workers[ProtoIDKeepAlive] = protocol.NewGetter()
	workers[ProtoIDGetUserInfo] = protocol.NewGetter()
	workers[ProtoIDVerification] = protocol.NewGetter()
	workers[ProtoIDGetDelayStatistics] = protocol.NewGetter()
}

var (
	ErrParameters    = errors.New("parameters missing or invalid")
	ErrInterrupted   = errors.New("process is interrupted")
	ErrChannelClosed = errors.New("channel is closed")
)

type OptionalInt32 struct {
	Value int32
}

type OptionalUInt64 struct {
	Value uint64
}

type OptionalDouble struct {
	Value float64
}

type OptionalBool struct {
	Value bool
}

// FutuAPI 是富途开放API的主要操作对象。
type FutuAPI struct {
	// 连接配置，通过方法设置，不设置默认为零值
	clientVer  int32
	clientID   string
	recvNotify bool
	encAlgo    common.PacketEncAlgo
	protoFmt   common.ProtoFmt

	// TCP连接，连接后设置
	connID uint64
	userID uint64
	// protocol
	proto *protocol.FutuProtocol

	// 发送心跳的定时器，连接后设置
	ticker *time.Ticker
	// 心跳定时器关闭信号通道
	done chan struct{}
}

// NewFutuAPI 创建API对象，并启动goroutine进行发送保活心跳.
func NewFutuAPI() *FutuAPI {
	return &FutuAPI{
		done: make(chan struct{}),
	}
}

// 设置调用接口信息, 非必调接口
func (api *FutuAPI) SetClientInfo(id string, ver int32) {
	api.clientID = id
	api.clientVer = ver
}

// 设置通讯协议 body 格式, 目前支持 Protobuf|Json 两种格式，默认 ProtoBuf, 非必调接口
func (api *FutuAPI) SetProtoFmt(fmt common.ProtoFmt) {
	api.protoFmt = fmt
}

// 获取连接 ID，连接初始化成功后才会有值
func (api *FutuAPI) ConnID() uint64 {
	return api.connID
}

func (api *FutuAPI) UserID() uint64 {
	return api.userID
}

func (api *FutuAPI) SetRecvNotify(recv bool) {
	api.recvNotify = recv
}

func (api *FutuAPI) SetEncAlgo(algo common.PacketEncAlgo) {
	api.encAlgo = algo
}

func (api *FutuAPI) packetID() *common.PacketID {
	return &common.PacketID{
		ConnID:   proto.Uint64(api.connID),
		SerialNo: proto.Uint32(api.proto.SerialNo()),
	}
}

// 连接FutuOpenD
func (api *FutuAPI) Connect(ctx context.Context, address string) error {
	proto, err := protocol.Connect(address, workers)
	if err != nil {
		return err
	}
	api.proto = proto
	resp, err := api.initConnect(ctx, api.clientVer, api.clientID, api.recvNotify, api.encAlgo, api.protoFmt, "go")
	if err != nil {
		return err
	}
	api.connID = resp.GetConnID()
	api.userID = resp.GetLoginUserID()
	if d := resp.GetKeepAliveInterval(); d > 0 {
		api.ticker = time.NewTicker(time.Second * time.Duration(d))
		go api.heartBeat(ctx)
	}
	return nil
}

// 关闭连接
func (api *FutuAPI) Close(ctx context.Context) error {
	if err := api.proto.Close(); err != nil {
		return err
	}
	close(api.done)
	return nil
}

func (api *FutuAPI) heartBeat(ctx context.Context) {
	for {
		select {
		case <-api.done:
			api.ticker.Stop()
			return
		case <-api.ticker.C:
			if _, err := api.keepAlive(ctx, time.Now().Unix()); err != nil {
				return
			}
		}
	}
}

// 初始化连接
func (api *FutuAPI) initConnect(ctx context.Context, clientVer int32, clientID string,
	recvNotify bool, encAlgo common.PacketEncAlgo, protoFmt common.ProtoFmt, lang string) (*initconnect.S2C, error) {

	if clientID == "" {
		return nil, ErrParameters
	}
	// 请求参数
	req := &initconnect.Request{
		C2S: &initconnect.C2S{
			ClientVer:     proto.Int32(clientVer),
			ClientID:      proto.String(clientID),
			RecvNotify:    proto.Bool(recvNotify),
			PacketEncAlgo: proto.Int32(int32(encAlgo)),
			PushProtoFmt:  proto.Int32(int32(protoFmt)),
		},
	}
	if lang != "" {
		req.C2S.ProgrammingLanguage = proto.String(lang)
	}

	// 发送请求，同步返回结果
	ch := make(chan *initconnect.Response)
	if err := api.proto.RegisterGet(ProtoIDInitConnect, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), protocol.Error(resp)
	}
}

// KeepAlive 保活心跳
func (api *FutuAPI) keepAlive(ctx context.Context, t int64) (int64, error) {
	// 请求参数
	req := &keepalive.Request{
		C2S: &keepalive.C2S{
			Time: proto.Int64(t),
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *keepalive.Response)
	if err := api.proto.RegisterGet(ProtoIDKeepAlive, req, protocol.NewProtobufChan(ch)); err != nil {
		return 0, err
	}
	select {
	case <-ctx.Done():
		return 0, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return 0, ErrChannelClosed
		}
		return resp.GetS2C().GetTime(), protocol.Error(resp)
	}
}

// 获取全局状态
func (api *FutuAPI) GetGlobalState(ctx context.Context) (*getglobalstate.S2C, error) {
	// 请求参数
	req := &getglobalstate.Request{
		C2S: &getglobalstate.C2S{
			UserID: proto.Uint64(api.UserID()),
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *getglobalstate.Response)
	if err := api.proto.RegisterGet(ProtoIDGetGlobalState, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), protocol.Error(resp)
	}
}

// 系统推送通知
func (api *FutuAPI) SysNotify(ctx context.Context) (<-chan *notify.Response, error) {
	ch := make(chan *notify.Response)
	if err := api.proto.RegisterUpdate(ProtoIDNotify, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	return ch, nil
}

// 获取延迟统计
func (api *FutuAPI) GetDelayStatistics(ctx context.Context, typeList []getdelaystatistics.DelayStatisticsType,
	pushStage getdelaystatistics.QotPushStage, segmentList []int32) (*getdelaystatistics.S2C, error) {

	if len(typeList) == 0 {
		return nil, ErrParameters
	}
	// 请求参数
	req := &getdelaystatistics.Request{
		C2S: &getdelaystatistics.C2S{
			TypeList: make([]int32, len(typeList)),
		},
	}
	for i, v := range typeList {
		req.C2S.TypeList[i] = int32(v)
	}
	if pushStage != getdelaystatistics.QotPushStage_QotPushStage_Unkonw {
		req.C2S.QotPushStage = proto.Int32(int32(pushStage))
	}
	if len(segmentList) != 0 {
		req.C2S.SegmentList = segmentList
	}

	// 发送请求，同步返回结果
	ch := make(chan *getdelaystatistics.Response)
	if err := api.proto.RegisterGet(ProtoIDGetDelayStatistics, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), protocol.Error(resp)
	}
}

func (api *FutuAPI) Verify(ctx context.Context, vType verification.VerificationType, op verification.VerificationOp,
	code string) error {

	if vType == verification.VerificationType_VerificationType_Unknow ||
		op == verification.VerificationOp_VerificationOp_Unknow ||
		(op == verification.VerificationOp_VerificationOp_InputAndLogin && code == "") {
		return ErrParameters
	}
	req := &verification.Request{
		C2S: &verification.C2S{
			Type: proto.Int32(int32(vType)),
			Op:   proto.Int32(int32(op)),
		},
	}
	if code != "" {
		req.C2S.Code = proto.String(code)
	}

	ch := make(chan *verification.Response)
	if err := api.proto.RegisterGet(ProtoIDVerification, req, protocol.NewProtobufChan(ch)); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return ErrChannelClosed
		}
		return protocol.Error(resp)
	}
}

// 获取用户信息, flag = 0返回全部信息
func (api *FutuAPI) GetUserInfo(ctx context.Context, flag getuserinfo.UserInfoField) (*getuserinfo.S2C, error) {
	// 请求参数
	req := &getuserinfo.Request{
		C2S: &getuserinfo.C2S{},
	}
	if flag != 0 {
		req.C2S.Flag = proto.Int32(int32(flag))
	}
	// 发送请求，同步返回结果
	ch := make(chan *getuserinfo.Response)
	if err := api.proto.RegisterGet(ProtoIDGetUserInfo, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), protocol.Error(resp)
	}
}
