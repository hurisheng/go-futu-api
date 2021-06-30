package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetbasicqot"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetBasicQot = 3004 //Qot_GetBasicQot	获取股票基本报价
)

// 获取股票基本行情
func (api *FutuAPI) GetStockQuote(ctx context.Context, securities []*Security) ([]*BasicQot, error) {
	// 请求参数
	req := qotgetbasicqot.Request{
		C2S: &qotgetbasicqot.C2S{
			SecurityList: securityList(securities).pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetbasicqot.ResponseChan)
	if err := api.get(ProtoIDQotGetBasicQot, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return basicQotListFromPB(resp.GetS2C().GetBasicQotList()), protocol.Error(resp)
	}
}
