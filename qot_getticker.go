package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetticker"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetTicker = 3010 //Qot_GetTicker	获取逐笔

func init() {
	workers[ProtoIDQotGetTicker] = protocol.NewGetter()
}

func (api *FutuAPI) GetRTTicker(ctx context.Context, security *qotcommon.Security, num int32) (*qotgetticker.S2C, error) {

	if security == nil {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetticker.Request{
		C2S: &qotgetticker.C2S{
			Security:  security,
			MaxRetNum: proto.Int32(num),
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetticker.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetTicker, req, protocol.NewProtobufChan(ch)); err != nil {
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
