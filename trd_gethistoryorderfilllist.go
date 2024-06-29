package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgethistoryorderfilllist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDTrdGetHistoryOrderFillList = 2222 //Trd_GetHistoryOrderFillList	获取历史成交列表

func init() {
	workers[ProtoIDTrdGetHistoryOrderFillList] = protocol.NewGetter()
}

// 查询历史成交
func (api *FutuAPI) GetHistoryDeal(ctx context.Context, header *trdcommon.TrdHeader, filter *trdcommon.TrdFilterConditions) ([]*trdcommon.OrderFill, error) {

	if header == nil || filter == nil {
		return nil, ErrParameters
	}
	req := &trdgethistoryorderfilllist.Request{
		C2S: &trdgethistoryorderfilllist.C2S{
			Header:           header,
			FilterConditions: filter,
		},
	}

	ch := make(chan *trdgethistoryorderfilllist.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetHistoryOrderFillList, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetOrderFillList(), protocol.Error(resp)
	}
}
