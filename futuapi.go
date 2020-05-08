package futuapi

import (
	"fmt"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/hurisheng/go-futu-api/protobuf/GetGlobalState"
	"github.com/hurisheng/go-futu-api/protobuf/InitConnect"
	"github.com/hurisheng/go-futu-api/protobuf/KeepAlive"
	"github.com/hurisheng/go-futu-api/protobuf/Notify"
)

// FutuAPI 是富途开放API的主要操作对象。
type FutuAPI struct {
	server *conn
	ticker *time.Ticker
	done   chan bool
}

// Config 为API配置信息
type Config struct {
	Address   string
	ClientVer int32
	ClientID  string
}

// NewFutuAPI 创建API对象，并启动goroutine进行发送保活心跳.
func NewFutuAPI(config *Config) (*FutuAPI, error) {
	// connect socket
	conn, err := newConn(config.Address)
	if err != nil {
		return nil, fmt.Errorf("connect to server error: %w", err)
	}
	api := &FutuAPI{
		server: conn,
		done:   make(chan bool),
	}
	// init connect
	ch, err := api.initConnect(&InitConnect.Request{
		C2S: &InitConnect.C2S{
			ClientVer: &config.ClientVer,
			ClientID:  &config.ClientID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("call InitConnect error: %w", err)
	}
	api.ticker = time.NewTicker(time.Second * time.Duration((<-ch).S2C.GetKeepAliveInterval()))
	go api.heartBeat()
	return api, nil
}

// Close 关闭连接.
func (api *FutuAPI) Close() error {
	api.done <- true
	if err := api.server.close(); err != nil {
		return fmt.Errorf("close server error: %w", err)
	}
	return nil
}

func (api *FutuAPI) heartBeat() {
	for {
		select {
		case <-api.done:
			api.ticker.Stop()
			return
		case <-api.ticker.C:
			now := time.Now().Unix()
			if _, err := api.keepAlive(&KeepAlive.Request{
				C2S: &KeepAlive.C2S{
					Time: &now,
				},
			}); err != nil {
				return
			}
		}
	}
}

func (api *FutuAPI) channel(out interface{}) (reflect.Value, reflect.Type, error) {
	// out必须为实现proto.Message接口的指针的channel
	// 通过reflect检查out的类型是否正确
	v, t := reflect.ValueOf(out), reflect.TypeOf(out)
	// out must be a channel type
	if t.Kind() != reflect.Chan {
		return reflect.ValueOf(nil), nil, fmt.Errorf("out is not channel type")
	}
	// it must be a channel of pointer to the response type which implements proto.Message interface
	p := t.Elem()
	if p.Kind() != reflect.Ptr || !p.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		return reflect.ValueOf(nil), nil, fmt.Errorf("out is not channel of pointer to type implements interface proto.Message")
	}
	return v, p.Elem(), nil
}

func (api *FutuAPI) send(protoID uint32, req proto.Message, out interface{}) error {
	// 传入的req为protobuf的数据类型，序列化后发送到服务器
	// 传入的out为实现proto.Message接口的结构类型指针的channel，根据实际的类型创建内存空间
	// 启动goroutine接收服务器返回的数据，并通过out channel发送
	buf, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request error: %w", err)
	}
	v, t, err := api.channel(out)
	if err != nil {
		return fmt.Errorf("parameter validation error: %w", err)
	}
	in, err := api.server.send(protoID, buf)
	if err != nil {
		return fmt.Errorf("send request error: %w", err)
	}
	resp := reflect.New(t)
	go func() {
		if err := proto.Unmarshal(<-in, resp.Interface().(proto.Message)); err != nil {
			v.Close()
			return
		}
		v.Send(resp)
	}()
	return nil
}

func (api *FutuAPI) subscribe(protoID uint32, out interface{}) error {
	// 传入的out为实现proto.Message接口的结构类型指针的channel，根据实际的类型创建内存空间
	// 在goroutine中不断接收服务器返回的数据，并通过out channel发送
	v, t, err := api.channel(out)
	if err != nil {
		return fmt.Errorf("parameter validation error: %w", err)
	}
	in, err := api.server.subscribe(protoID)
	if err != nil {
		return fmt.Errorf("subscribe error: %w", err)
	}
	go func() {
		for buf := range in {
			resp := reflect.New(t)
			if err := proto.Unmarshal(buf, resp.Interface().(proto.Message)); err != nil {
				break
			}
			v.Send(resp)
		}
		v.Close()
	}()
	return nil
}

// InitConnect 初始化连接
func (api *FutuAPI) initConnect(req *InitConnect.Request) (<-chan *InitConnect.Response, error) {
	out := make(chan *InitConnect.Response)
	if err := api.send(ProtoIDInitConnect, req, out); err != nil {
		return nil, fmt.Errorf("InitConnect error: %w", err)
	}
	return out, nil
}

// KeepAlive 保活心跳
func (api *FutuAPI) keepAlive(req *KeepAlive.Request) (<-chan *KeepAlive.Response, error) {
	out := make(chan *KeepAlive.Response)
	if err := api.send(ProtoIDKeepAlive, req, out); err != nil {
		return nil, fmt.Errorf("KeepAlive error: %w", err)
	}
	return out, nil
}

// GetGlobalState 获取全局状态
func (api *FutuAPI) GetGlobalState(req *GetGlobalState.Request) (<-chan *GetGlobalState.Response, error) {
	out := make(chan *GetGlobalState.Response)
	if err := api.send(ProtoIDGetGlobalState, req, out); err != nil {
		return nil, fmt.Errorf("GetGlobalState error: %w", err)
	}
	return out, nil
}

// Notify 系统推送通知
func (api *FutuAPI) Notify() (<-chan *Notify.Response, error) {
	out := make(chan *Notify.Response)
	if err := api.subscribe(ProtoIDNotify, out); err != nil {
		return nil, fmt.Errorf("Notify error: %w", err)
	}
	return out, nil
}
