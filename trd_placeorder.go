package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdplaceorder"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdPlaceOrder = 2202 //Trd_PlaceOrder	下单

func init() {
	workers[ProtoIDTrdPlaceOrder] = protocol.NewGetter()
}

// 下单
func (api *FutuAPI) PlaceOrder(ctx context.Context, header *trdcommon.TrdHeader, trdSide trdcommon.TrdSide, orderType trdcommon.OrderType, code string, qty float64,
	price *OptionalDouble, adjustPrice *OptionalBool, sideAndLimit *OptionalDouble, secMarket trdcommon.TrdSecMarket, remark string, timeInForce trdcommon.TimeInForce,
	fillOutsideRTH *OptionalBool, auxPrice *OptionalDouble, trailType trdcommon.TrailType, trailValue *OptionalDouble, trailSpread *OptionalDouble) (*trdplaceorder.S2C, error) {

	if header == nil ||
		trdSide == trdcommon.TrdSide_TrdSide_Unknown ||
		orderType == trdcommon.OrderType_OrderType_Unknown ||
		code == "" {
		return nil, ErrParameters
	}
	req := &trdplaceorder.Request{
		C2S: &trdplaceorder.C2S{
			PacketID:  api.packetID(),
			Header:    header,
			TrdSide:   proto.Int32(int32(trdSide)),
			OrderType: proto.Int32(int32(orderType)),
			Code:      proto.String(code),
			Qty:       proto.Float64(qty),
		},
	}
	if price != nil {
		req.C2S.Price = proto.Float64(price.Value)
	}
	if adjustPrice != nil {
		req.C2S.AdjustPrice = proto.Bool(adjustPrice.Value)
	}
	if sideAndLimit != nil {
		req.C2S.AdjustSideAndLimit = proto.Float64(sideAndLimit.Value)
	}
	if secMarket != trdcommon.TrdSecMarket_TrdSecMarket_Unknown {
		req.C2S.SecMarket = proto.Int32(int32(secMarket))
	}
	if remark != "" {
		req.C2S.Remark = proto.String(remark)
	}
	if timeInForce == trdcommon.TimeInForce_TimeInForce_DAY || timeInForce == trdcommon.TimeInForce_TimeInForce_GTC {
		req.C2S.TimeInForce = proto.Int32(int32(timeInForce))
	}
	if fillOutsideRTH != nil {
		req.C2S.FillOutsideRTH = proto.Bool(fillOutsideRTH.Value)
	}
	if auxPrice != nil {
		req.C2S.AuxPrice = proto.Float64(auxPrice.Value)
	}
	if trailType != trdcommon.TrailType_TrailType_Unknown {
		req.C2S.TrailType = proto.Int32(int32(trailType))
	}
	if trailValue != nil {
		req.C2S.TrailValue = proto.Float64(trailValue.Value)
	}
	if trailSpread != nil {
		req.C2S.TrailSpread = proto.Float64(trailSpread.Value)
	}

	ch := make(chan *trdplaceorder.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdPlaceOrder, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), protocol.Error(resp)
	}
}
