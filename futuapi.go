package futuapi

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/hurisheng/go-futu-api/pb/common"
	"github.com/hurisheng/go-futu-api/pb/getglobalstate"
	"github.com/hurisheng/go-futu-api/pb/initconnect"
	"github.com/hurisheng/go-futu-api/pb/keepalive"
	"github.com/hurisheng/go-futu-api/pb/notify"
	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/protocol"
	"github.com/hurisheng/go-futu-api/tcp"
)

var (
	ErrInterrupted   = errors.New("process is interrupted")
	ErrChannelClosed = errors.New("channel is closed")
)

// FutuAPI 是富途开放API的主要操作对象。
type FutuAPI struct {
	// 连接配置，通过方法设置，不设置默认为零值
	clientVer  int32
	clientID   string
	recvNotify bool
	encAlgo    common.PacketEncAlgo
	protoFmt   common.ProtoFmt

	// TCP连接，连接后设置
	conn   *tcp.Conn
	connID uint64
	userID uint64
	// 数据接收注册表
	reg *protocol.Registry

	serial uint32
	mu     sync.Mutex
	// 发送心跳的定时器，连接后设置
	ticker *time.Ticker
	// 心跳定时器关闭信号通道
	done chan struct{}
}

// NewFutuAPI 创建API对象，并启动goroutine进行发送保活心跳.
func NewFutuAPI() *FutuAPI {
	return &FutuAPI{
		reg:    protocol.NewRegistry(),
		done:   make(chan struct{}),
		serial: 1,
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

func (api *FutuAPI) serialNo() uint32 {
	// 递增serial
	api.mu.Lock()
	defer api.mu.Unlock()
	api.serial++
	return api.serial
}

// 连接FutuOpenD
func (api *FutuAPI) Connect(ctx context.Context, address string) error {
	conn, err := tcp.Dial(address, protocol.NewDecoder(api.reg))
	if err != nil {
		return err
	}
	api.conn = conn
	resp, err := api.initConnect(ctx, api.clientVer, api.clientID, api.recvNotify, api.encAlgo, api.protoFmt, "golang")
	if err != nil {
		return err
	}
	api.connID = resp.ConnID
	api.userID = resp.LoginUserID
	if d := resp.KeepAliveInterval; d > 0 {
		api.ticker = time.NewTicker(time.Second * time.Duration(d))
		go api.heartBeat(ctx)
	}
	return nil
}

// 关闭连接
func (api *FutuAPI) Close(ctx context.Context) error {
	if err := api.conn.Close(); err != nil {
		return err
	}
	close(api.done)
	api.reg.Close()
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

func (api *FutuAPI) get(proto uint32, req proto.Message, out protocol.RespChan) error {
	// 获取serial
	se := api.serialNo()
	// 在registry注册get channel
	if err := api.reg.AddGetChan(proto, se, out); err != nil {
		return err
	}
	// 向服务器发送req
	if err := api.conn.Write(protocol.NewEncoder(proto, se, req)); err != nil {
		if err := api.reg.RemoveChan(proto, se); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (api *FutuAPI) update(proto uint32, out protocol.RespChan) error {
	// 在registry注册update channel
	if err := api.reg.AddUpdateChan(proto, out); err != nil {
		return err
	}
	return nil
}

const (
	ProtoIDInitConnect    = 1001 //InitConnect	初始化连接
	ProtoIDGetGlobalState = 1002 //GetGlobalState	获取全局状态
	ProtoIDNotify         = 1003 //Notify	系统通知推送
	ProtoIDKeepAlive      = 1004 //KeepAlive	保活心跳
)

// 初始化连接
func (api *FutuAPI) initConnect(ctx context.Context, clientVer int32, clientID string, recvNotify bool,
	encAlgo common.PacketEncAlgo, protoFmt common.ProtoFmt, lang string) (*initConnectResp, error) {
	// 请求参数
	req := initconnect.Request{C2S: &initconnect.C2S{
		ClientVer:           &clientVer,
		ClientID:            &clientID,
		RecvNotify:          &recvNotify,
		PacketEncAlgo:       (*int32)(&encAlgo),
		PushProtoFmt:        (*int32)(&protoFmt),
		ProgrammingLanguage: &lang,
	}}
	// 发送请求，同步返回结果
	ch := make(initconnect.ResponseChan)
	if err := api.get(ProtoIDInitConnect, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return initConnectRespFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type initConnectResp struct {
	ServerVer         int32  //FutuOpenD 的版本号
	LoginUserID       uint64 //FutuOpenD 登陆的牛牛用户 ID
	ConnID            uint64 //此连接的连接 ID，连接的唯一标识
	ConnAESKey        string //此连接后续 AES 加密通信的 Key，固定为16字节长字符串
	KeepAliveInterval int32  //心跳保活间隔
	AESCBCiv          string //AES 加密通信 CBC 加密模式的 iv，固定为16字节长字符串
}

func initConnectRespFromPB(pb *initconnect.S2C) *initConnectResp {
	if pb == nil {
		return nil
	}
	return &initConnectResp{
		ServerVer:         pb.GetServerVer(),
		LoginUserID:       pb.GetLoginUserID(),
		ConnID:            pb.GetConnID(),
		ConnAESKey:        pb.GetConnAESKey(),
		KeepAliveInterval: pb.GetKeepAliveInterval(),
		AESCBCiv:          pb.GetAesCBCiv(),
	}
}

// KeepAlive 保活心跳
func (api *FutuAPI) keepAlive(ctx context.Context, t int64) (int64, error) {
	// 请求参数
	req := keepalive.Request{C2S: &keepalive.C2S{
		Time: &t,
	}}
	// 发送请求，同步返回结果
	ch := make(keepalive.ResponseChan)
	if err := api.get(ProtoIDKeepAlive, &req, ch); err != nil {
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
func (api *FutuAPI) GetGlobalState(ctx context.Context) (*GlobalState, error) {
	// 请求参数
	req := getglobalstate.Request{C2S: &getglobalstate.C2S{
		UserID: &api.userID,
	}}
	// 发送请求，同步返回结果
	ch := make(getglobalstate.ResponseChan)
	if err := api.get(ProtoIDGetGlobalState, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return globalStateFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type GlobalState struct {
	MarketHK       qotcommon.QotMarketState
	MarketUS       qotcommon.QotMarketState
	MarketSH       qotcommon.QotMarketState
	MarketSZ       qotcommon.QotMarketState
	MarketHKFuture qotcommon.QotMarketState
	MarketUSFuture qotcommon.QotMarketState
	QotLogined     bool
	TrdLogined     bool
	ServerVer      int32
	ServerBuildNo  int32
	Time           int64
	LocalTime      float64
	ProgramStatus  *ProgramStatus
	QotSvrIpAddr   string
	TrdSvrIpAddr   string
	ConnID         uint64
}

func globalStateFromPB(resp *getglobalstate.S2C) *GlobalState {
	if resp == nil {
		return nil
	}
	return &GlobalState{
		MarketHK:       qotcommon.QotMarketState(resp.GetMarketHK()),
		MarketUS:       qotcommon.QotMarketState(resp.GetMarketUS()),
		MarketSH:       qotcommon.QotMarketState(resp.GetMarketSH()),
		MarketSZ:       qotcommon.QotMarketState(resp.GetMarketSZ()),
		MarketHKFuture: qotcommon.QotMarketState(resp.GetMarketHKFuture()),
		MarketUSFuture: qotcommon.QotMarketState(resp.GetMarketUSFuture()),
		QotLogined:     resp.GetQotLogined(),
		TrdLogined:     resp.GetTrdLogined(),
		ServerVer:      resp.GetServerVer(),
		ServerBuildNo:  resp.GetServerBuildNo(),
		Time:           resp.GetTime(),
		LocalTime:      resp.GetLocalTime(),
		ProgramStatus:  programStatusFromPB(resp.GetProgramStatus()),
		QotSvrIpAddr:   resp.GetQotSvrIpAddr(),
		TrdSvrIpAddr:   resp.GetTrdSvrIpAddr(),
		ConnID:         resp.GetConnID(),
	}
}

// 系统推送通知
func (api *FutuAPI) SysNotify(ctx context.Context) (<-chan *SysNotifyResp, error) {
	ch := make(notifyChan)
	if err := api.update(ProtoIDNotify, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type SysNotifyResp struct {
	Notification *Notification
	Err          error
}

type notifyChan chan *SysNotifyResp

var _ protocol.RespChan = make(notifyChan)

func (ch notifyChan) Send(b []byte) error {
	var resp notify.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &SysNotifyResp{
		Notification: notificationFromPB(resp.GetS2C()),
		Err:          protocol.Error(&resp),
	}
	return nil
}

func (ch notifyChan) Close() {
	close(ch)
}

type Notification struct {
	Type          notify.NotifyType    //*通知类型
	Event         *GtwEvent            //事件通息
	ProgramStatus *NotifyProgramStatus //程序状态
	ConnectStatus *ConnectStatus       //连接状态
	QotRight      *QotRight            //行情权限
	APILevel      *APILevel            //用户等级，已在2.10版本之后废弃
	APIQuota      *APIQuota            //API 额度
}

func notificationFromPB(pb *notify.S2C) *Notification {
	if pb == nil {
		return nil
	}
	return &Notification{
		Type:          notify.NotifyType(pb.GetType()),
		Event:         gtwEventFromPB(pb.GetEvent()),
		ProgramStatus: notifyProgramStatusFromPB(pb.GetProgramStatus()),
		ConnectStatus: connectStatusFromPB(pb.GetConnectStatus()),
		QotRight:      qotRightFromPB(pb.GetQotRight()),
		APILevel:      apiLevelFromPB(pb.GetApiLevel()),
		APIQuota:      apiQuotaFromPB(pb.GetApiQuota()),
	}
}

type GtwEvent struct {
	EventType notify.GtwEventType //*GtwEventType,事件类型
	Desc      string              //*事件描述
}

func gtwEventFromPB(pb *notify.GtwEvent) *GtwEvent {
	if pb == nil {
		return nil
	}
	return &GtwEvent{
		EventType: notify.GtwEventType(pb.GetEventType()),
		Desc:      pb.GetDesc(),
	}
}

type NotifyProgramStatus struct {
	ProgramStatus *ProgramStatus //*当前程序状态
}

func notifyProgramStatusFromPB(pb *notify.ProgramStatus) *NotifyProgramStatus {
	if pb == nil {
		return nil
	}
	return &NotifyProgramStatus{
		ProgramStatus: programStatusFromPB(pb.GetProgramStatus()),
	}
}

type ConnectStatus struct {
	QotLogined bool //*是否登陆行情服务器
	TrdLogined bool //*是否登陆交易服务器
}

func connectStatusFromPB(pb *notify.ConnectStatus) *ConnectStatus {
	if pb == nil {
		return nil
	}
	return &ConnectStatus{
		QotLogined: pb.GetQotLogined(),
		TrdLogined: pb.GetTrdLogined(),
	}
}

type QotRight struct {
	HKQotRight          qotcommon.QotRight //*港股行情权限, Qot_Common.QotRight
	USQotRight          qotcommon.QotRight //*美股行情权限, Qot_Common.QotRight
	CNQotRight          qotcommon.QotRight //*A股行情权限, Qot_Common.QotRight
	HKOptionQotRight    qotcommon.QotRight //港股期权行情权限, Qot_Common.QotRight
	HasUSOptionQotRight bool               //是否有美股期权行情权限
	HKFutureQotRight    qotcommon.QotRight //港股期货行情权限, Qot_Common.QotRight
	USFutureQotRight    qotcommon.QotRight //美股期货行情权限, Qot_Common.QotRight
	USOptionQotRight    qotcommon.QotRight //美股期货行情权限, Qot_Common.QotRight
}

func qotRightFromPB(pb *notify.QotRight) *QotRight {
	if pb == nil {
		return nil
	}
	return &QotRight{
		HKQotRight:          qotcommon.QotRight(pb.GetHkQotRight()),
		USQotRight:          qotcommon.QotRight(pb.GetUsQotRight()),
		CNQotRight:          qotcommon.QotRight(pb.GetCnQotRight()),
		HKOptionQotRight:    qotcommon.QotRight(pb.GetHkOptionQotRight()),
		HasUSOptionQotRight: pb.GetHasUSOptionQotRight(),
		HKFutureQotRight:    qotcommon.QotRight(pb.GetHkFutureQotRight()),
		USFutureQotRight:    qotcommon.QotRight(pb.GetUsFutureQotRight()),
		USOptionQotRight:    qotcommon.QotRight(pb.GetUsOptionQotRight()),
	}
}

type APILevel struct {
	APILevel string //*api用户等级描述，已在2.10版本之后废弃
}

func apiLevelFromPB(pb *notify.APILevel) *APILevel {
	if pb == nil {
		return nil
	}
	return &APILevel{
		APILevel: pb.GetApiLevel(),
	}
}

type APIQuota struct {
	SubQuota       int32 //*订阅额度
	HistoryKLQuota int32 //*历史K线额度
}

func apiQuotaFromPB(pb *notify.APIQuota) *APIQuota {
	if pb == nil {
		return nil
	}
	return &APIQuota{
		SubQuota:       pb.GetSubQuota(),
		HistoryKLQuota: pb.GetHistoryKLQuota(),
	}
}
