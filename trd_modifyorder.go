package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdmodifyorder"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdModifyOrder = 2205 //Trd_ModifyOrder

func init() {
	workers[ProtoIDTrdModifyOrder] = protocol.NewGetter()
}

// 改单撤单
func (api *FutuAPI) ModifyOrder(ctx context.Context, header *trdcommon.TrdHeader, orderID uint64, op trdcommon.ModifyOrderOp, forAll bool, 
	trdMarket trdcommon.TrdMarket, qty *OptionalDouble, price *OptionalDouble,
	adjust *OptionalBool, sideAndLimit *OptionalDouble, auxPrice *OptionalDouble,
	trailType trdcommon.TrailType, trailValue *OptionalDouble, trailSpread *OptionalDouble, orderIDEx string) (*trdmodifyorder.S2C, error) {

	if header == nil || op == trdcommon.ModifyOrderOp_ModifyOrderOp_Unknown {
		return nil, ErrParameters
	}
	req := &trdmodifyorder.Request{
		C2S: &trdmodifyorder.C2S{
			PacketID:      api.packetID(),
			Header:        header,
			OrderID:       proto.Uint64(orderID),
			ModifyOrderOp: proto.Int32(int32(op)),
			ForAll:        proto.Bool(forAll),
		},
	}
	if trdMarket != trdcommon.TrdMarket_TrdMarket_Unknown {
		req.C2S.TrdMarket = proto.Int32(int32(trdMarket))
	}
	if qty != nil {
		req.C2S.Qty = proto.Float64(qty.Value)
	}
	if price != nil {
		req.C2S.Price = proto.Float64(price.Value)
	}
	if adjust != nil {
		req.C2S.AdjustPrice = proto.Bool(adjust.Value)
	}
	if sideAndLimit != nil {
		req.C2S.AdjustSideAndLimit = proto.Float64(sideAndLimit.Value)
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
	if orderIDEx != "" {
		req.C2S.OrderIDEx = proto.String(orderIDEx)
	}

	ch := make(chan *trdmodifyorder.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdModifyOrder, req, protocol.NewProtobufChan(ch)); err != nil {
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
