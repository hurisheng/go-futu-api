package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetcapitaldistribution"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetCapitalDistribution = 3212 //Qot_GetCapitalDistribution 获取资金分布

func init() {
	workers[ProtoIDQotGetCapitalDistribution] = protocol.NewGetter()
}

// 获取资金分布
func (api *FutuAPI) GetCapitalDistribution(ctx context.Context, security *qotcommon.Security) (*qotgetcapitaldistribution.S2C, error) {

	if security == nil {
		return nil, ErrParameters
	}
	// 请求数据
	req := &qotgetcapitaldistribution.Request{
		C2S: &qotgetcapitaldistribution.C2S{
			Security: security,
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetcapitaldistribution.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetCapitalDistribution, req, protocol.NewProtobufChan(ch)); err != nil {
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
