package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetmarketstate"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetMarketState = 3223 //Qot_GetMarketState	获取指定品种的市场状态

func init() {
	workers[ProtoIDQotGetMarketState] = protocol.NewGetter()
}

// 获取标的市场状态
func (api *FutuAPI) GetMarketState(ctx context.Context, securities []*qotcommon.Security) ([]*qotgetmarketstate.MarketInfo, error) {

	if len(securities) == 0 {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetmarketstate.Request{
		C2S: &qotgetmarketstate.C2S{
			SecurityList: securities,
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetmarketstate.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetMarketState, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetMarketInfoList(), protocol.Error(resp)
	}
}
