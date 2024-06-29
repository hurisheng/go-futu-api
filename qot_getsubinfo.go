package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetsubinfo"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetSubInfo = 3003 //Qot_GetSubInfo	获取订阅信息

func init() {
	workers[ProtoIDQotGetSubInfo] = protocol.NewGetter()
}

// 获取订阅信息
func (api *FutuAPI) QuerySubscription(ctx context.Context, isAll bool) (*qotgetsubinfo.S2C, error) {
	// 请求参数
	req := &qotgetsubinfo.Request{
		C2S: &qotgetsubinfo.C2S{
			IsReqAllConn: proto.Bool(isAll),
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetsubinfo.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetSubInfo, req, protocol.NewProtobufChan(ch)); err != nil {
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
