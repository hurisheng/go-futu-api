package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetrt"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetRT = 3008 //Qot_GetRT	获取分时

func init() {
	workers[ProtoIDQotGetRT] = protocol.NewGetter()
}

// 获取实时分时
func (api *FutuAPI) GetRTData(ctx context.Context, security *qotcommon.Security) (*qotgetrt.S2C, error) {
	if security == nil {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetrt.Request{
		C2S: &qotgetrt.C2S{
			Security: security,
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetrt.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetRT, req, protocol.NewProtobufChan(ch)); err != nil {
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
