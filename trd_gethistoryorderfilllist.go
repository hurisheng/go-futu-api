package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdgethistoryorderfilllist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdGetHistoryOrderFillList = 2222 //Trd_GetHistoryOrderFillList	获取历史成交列表
)

// 查询历史成交
func (api *FutuAPI) GetHistoryDeal(ctx context.Context, header *TrdHeader, filter *TrdFilterConditions) ([]*OrderFill, error) {
	req := trdgethistoryorderfilllist.Request{
		C2S: &trdgethistoryorderfilllist.C2S{
			Header:           header.pb(),
			FilterConditions: filter.pb(),
		},
	}
	ch := make(trdgethistoryorderfilllist.ResponseChan)
	if err := api.get(ProtoIDTrdGetHistoryOrderFillList, &req, ch); err != nil {
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
