package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdgetorderfilllist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdGetOrderFillList = 2211 //Trd_GetOrderFillList	获取成交列表
)

// 查询当日成交
func (api *FutuAPI) GetDealList(ctx context.Context, header *TrdHeader, filter *TrdFilterConditions, refresh bool) ([]*OrderFill, error) {
	req := trdgetorderfilllist.Request{
		C2S: &trdgetorderfilllist.C2S{
			Header:           header.pb(),
			FilterConditions: filter.pb(),
			RefreshCache:     &refresh,
		},
	}
	ch := make(trdgetorderfilllist.ResponseChan)
	if err := api.get(ProtoIDTrdGetOrderFillList, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return orderFillListFromPB(resp.GetS2C().GetOrderFillList()), protocol.Error(resp)
	}
}
