package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetbasicqot"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetBasicQot = 3004 //Qot_GetBasicQot	获取股票基本报价

func init() {
	workers[ProtoIDQotGetBasicQot] = protocol.NewGetter()
}

// 获取股票基本行情
func (api *FutuAPI) GetStockQuote(ctx context.Context, securities []*qotcommon.Security) ([]*qotcommon.BasicQot, error) {
	if len(securities) == 0 {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetbasicqot.Request{
		C2S: &qotgetbasicqot.C2S{
			SecurityList: securities,
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetbasicqot.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetBasicQot, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetBasicQotList(), protocol.Error(resp)
	}
}
