package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgethistoryorderlist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDTrdGetHistoryOrderList = 2221 //Trd_GetHistoryOrderList	获取历史订单列表

func init() {
	workers[ProtoIDTrdGetHistoryOrderList] = protocol.NewGetter()
}

// 查询历史订单
func (api *FutuAPI) GetHistoryOrderList(ctx context.Context, header *trdcommon.TrdHeader, filter *trdcommon.TrdFilterConditions,
	status []trdcommon.OrderStatus) ([]*trdcommon.Order, error) {

	if header == nil || filter == nil {
		return nil, ErrParameters
	}
	req := &trdgethistoryorderlist.Request{
		C2S: &trdgethistoryorderlist.C2S{
			Header:           header,
			FilterConditions: filter,
		},
	}
	if len(status) != 0 {
		req.C2S.FilterStatusList = make([]int32, len(status))
		for i, v := range status {
			req.C2S.FilterStatusList[i] = int32(v)
		}
	}

	ch := make(chan *trdgethistoryorderlist.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetHistoryOrderList, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetOrderList(), protocol.Error(resp)
	}
}
