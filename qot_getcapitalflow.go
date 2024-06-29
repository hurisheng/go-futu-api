package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetcapitalflow"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetCapitalFlow = 3211 //Qot_GetCapitalFlow	获取资金流向

func init() {
	workers[ProtoIDQotGetCapitalFlow] = protocol.NewGetter()
}

// 获取资金流向
func (api *FutuAPI) GetCapitalFlow(ctx context.Context, security *qotcommon.Security,
	periodType qotcommon.PeriodType, begin string, end string) (*qotgetcapitalflow.S2C, error) {

	if security == nil {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetcapitalflow.Request{
		C2S: &qotgetcapitalflow.C2S{
			Security: security,
		},
	}
	if periodType != qotcommon.PeriodType_PeriodType_Unknown {
		req.C2S.PeriodType = proto.Int32(int32(periodType))
	}
	if begin != "" {
		req.C2S.BeginTime = proto.String(begin)
	}
	if end != "" {
		req.C2S.EndTime = proto.String(end)
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetcapitalflow.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetCapitalFlow, req, protocol.NewProtobufChan(ch)); err != nil {
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
