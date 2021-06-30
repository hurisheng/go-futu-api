package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdmodifyorder"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdModifyOrder = 2205 //Trd_ModifyOrder
)

// 改单撤单
func (api *FutuAPI) ModifyOrder(ctx context.Context, header *TrdHeader, all bool, op trdcommon.ModifyOrderOp, orderID uint64,
	qty float64, price float64, adjustPrice bool, sideAndLimit float64) (uint64, error) {
	req := trdmodifyorder.Request{
		C2S: &trdmodifyorder.C2S{
			PacketID:           api.packetID().pb(),
			Header:             header.pb(),
			OrderID:            &orderID,
			ModifyOrderOp:      (*int32)(&op),
			ForAll:             &all,
			Qty:                &qty,
			Price:              &price,
			AdjustPrice:        &adjustPrice,
			AdjustSideAndLimit: &sideAndLimit,
		},
	}
	ch := make(trdmodifyorder.ResponseChan)
	if err := api.get(ProtoIDTrdModifyOrder, &req, ch); err != nil {
		return 0, err
	}
	select {
	case <-ctx.Done():
		return 0, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return 0, ErrChannelClosed
		}
		return resp.GetS2C().GetOrderID(), protocol.Error(resp)
	}
}
