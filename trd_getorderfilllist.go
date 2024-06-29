package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetorderfilllist"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdGetOrderFillList = 2211 //Trd_GetOrderFillList	获取成交列表

func init() {
	workers[ProtoIDTrdGetOrderFillList] = protocol.NewGetter()
}

// 查询当日成交
func (api *FutuAPI) GetDealList(ctx context.Context, header *trdcommon.TrdHeader,
	filter *trdcommon.TrdFilterConditions, refresh *OptionalBool) ([]*trdcommon.OrderFill, error) {

	if header == nil {
		return nil, ErrParameters
	}
	req := &trdgetorderfilllist.Request{
		C2S: &trdgetorderfilllist.C2S{
			Header:           header,
			FilterConditions: filter,
		},
	}
	if refresh != nil {
		req.C2S.RefreshCache = proto.Bool(refresh.Value)
	}

	ch := make(chan *trdgetorderfilllist.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetOrderFillList, req, protocol.NewProtobufChan(ch)); err != nil {
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
