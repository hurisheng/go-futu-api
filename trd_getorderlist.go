package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetorderlist"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdGetOrderList = 2201 //Trd_GetOrderList	获取订单列表

func init() {
	workers[ProtoIDTrdGetOrderList] = protocol.NewGetter()
}

// 查询今日订单
func (api *FutuAPI) GetOrderList(ctx context.Context, header *trdcommon.TrdHeader,
	filter *trdcommon.TrdFilterConditions, status []trdcommon.OrderStatus, refresh *OptionalBool) ([]*trdcommon.Order, error) {

	if header == nil {
		return nil, ErrParameters
	}
	req := &trdgetorderlist.Request{
		C2S: &trdgetorderlist.C2S{
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
	if refresh != nil {
		req.C2S.RefreshCache = proto.Bool(refresh.Value)
	}

	ch := make(chan *trdgetorderlist.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetOrderList, req, protocol.NewProtobufChan(ch)); err != nil {
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
