package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgethistoryorderlist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdGetHistoryOrderList = 2221 //Trd_GetHistoryOrderList	获取历史订单列表
)

// 查询历史订单
func (api *FutuAPI) GetHistoryOrderList(ctx context.Context, header *TrdHeader, filter *TrdFilterConditions, status []trdcommon.OrderStatus) ([]*Order, error) {
	req := trdgethistoryorderlist.Request{
		C2S: &trdgethistoryorderlist.C2S{
			Header:           header.pb(),
			FilterConditions: filter.pb(),
			FilterStatusList: orderStatusList(status).pb(),
		},
	}
	ch := make(trdgethistoryorderlist.ResponseChan)
	if err := api.get(ProtoIDTrdGetHistoryOrderList, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return orderListFromPB(resp.GetS2C().GetOrderList()), protocol.Error(resp)
	}
}
