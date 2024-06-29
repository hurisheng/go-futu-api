package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetbroker"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetBroker = 3014 //Qot_GetBroker	获取经纪队列

func init() {
	workers[ProtoIDQotGetBroker] = protocol.NewGetter()
}

// 获取实时经纪队列
func (api *FutuAPI) GetBrokerQueue(ctx context.Context, security *qotcommon.Security) (*qotgetbroker.S2C, error) {

	if security == nil {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetbroker.Request{
		C2S: &qotgetbroker.C2S{
			Security: security,
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetbroker.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetBroker, req, protocol.NewProtobufChan(ch)); err != nil {
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
