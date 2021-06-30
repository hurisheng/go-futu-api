package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdplaceorder"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdPlaceOrder = 2202 //Trd_PlaceOrder	下单
)

// 下单
func (api *FutuAPI) PlaceOrder(ctx context.Context, header *TrdHeader, trdSide trdcommon.TrdSide, orderType trdcommon.OrderType, code string, qty float64, price float64,
	adjustPrice bool, sideAndLimit float64, secMarket trdcommon.TrdSecMarket, remark string, timeInForce trdcommon.TimeInForce, fillOutsideRTH bool) (uint64, error) {
	req := trdplaceorder.Request{
		C2S: &trdplaceorder.C2S{
			PacketID:           api.packetID().pb(),
			Header:             header.pb(),
			TrdSide:            (*int32)(&trdSide),
			OrderType:          (*int32)(&orderType),
			Code:               &code,
			Qty:                &qty,
			AdjustPrice:        &adjustPrice,
			AdjustSideAndLimit: &sideAndLimit,
			Remark:             &remark,
			FillOutsideRTH:     &fillOutsideRTH,
		},
	}
	if secMarket != 0 {
		req.C2S.SecMarket = (*int32)(&secMarket)
	}
	if timeInForce != 0 {
		req.C2S.TimeInForce = (*int32)(&timeInForce)
	}
	ch := make(trdplaceorder.ResponseChan)
	if err := api.get(ProtoIDTrdPlaceOrder, &req, ch); err != nil {
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
