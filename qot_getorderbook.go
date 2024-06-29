package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetorderbook"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetOrderBook = 3012 //Qot_GetOrderBook	获取买卖盘

func init() {
	workers[ProtoIDQotGetOrderBook] = protocol.NewGetter()
}

// 获取实时摆盘
func (api *FutuAPI) GetOrderBook(ctx context.Context, security *qotcommon.Security, num int32) (*qotgetorderbook.S2C, error) {

	if security == nil {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetorderbook.Request{
		C2S: &qotgetorderbook.C2S{
			Security: security,
			Num:      proto.Int32(num),
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetorderbook.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetOrderBook, req, protocol.NewProtobufChan(ch)); err != nil {
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
