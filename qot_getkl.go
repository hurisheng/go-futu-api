package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetkl"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetKL = 3006 //Qot_GetKL	获取 K 线

func init() {
	workers[ProtoIDQotGetKL] = protocol.NewGetter()
}

// 获取实时 K 线
func (api *FutuAPI) GetCurKLine(ctx context.Context, security *qotcommon.Security, num int32, rehabType qotcommon.RehabType, klType qotcommon.KLType) (*qotgetkl.S2C, error) {
	if security == nil ||
		rehabType == qotcommon.RehabType_RehabType_None ||
		klType == qotcommon.KLType_KLType_Unknown {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetkl.Request{
		C2S: &qotgetkl.C2S{
			Security:  security,
			ReqNum:    proto.Int32(num),
			RehabType: proto.Int32(int32(rehabType)),
			KlType:    proto.Int32(int32(klType)),
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetkl.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetKL, req, protocol.NewProtobufChan(ch)); err != nil {
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
