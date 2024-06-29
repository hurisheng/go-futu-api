package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetmaxtrdqtys"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdGetMaxTrdQtys = 2111 //Trd_GetMaxTrdQtys	获取最大交易数量

func init() {
	workers[ProtoIDTrdGetMaxTrdQtys] = protocol.NewGetter()
}

// 查询最大可买可卖
func (api *FutuAPI) GetMaxTrdQtys(ctx context.Context, header *trdcommon.TrdHeader, orderType trdcommon.OrderType, code string, price float64,
	orderID *OptionalUInt64, adjust *OptionalBool, sideAndLimit *OptionalDouble, secMarket trdcommon.TrdSecMarket, orderIDEx string) (*trdcommon.MaxTrdQtys, error) {
	// required parameters should not be invalid
	if header == nil || orderType == trdcommon.OrderType_OrderType_Unknown || code == "" {
		return nil, ErrParameters
	}
	// request information
	req := &trdgetmaxtrdqtys.Request{
		C2S: &trdgetmaxtrdqtys.C2S{
			Header:    header,
			OrderType: proto.Int32(int32(orderType)),
			Code:      proto.String(code),
			Price:     proto.Float64(price),
		},
	}
	// optional parameters
	if orderID != nil {
		req.C2S.OrderID = proto.Uint64(orderID.Value)
	}
	if adjust != nil {
		req.C2S.AdjustPrice = proto.Bool(adjust.Value)
	}
	if sideAndLimit != nil {
		req.C2S.AdjustSideAndLimit = proto.Float64(sideAndLimit.Value)
	}
	if secMarket != trdcommon.TrdSecMarket_TrdSecMarket_Unknown {
		req.C2S.SecMarket = proto.Int32(int32(secMarket))
	}
	if orderIDEx != "" {
		req.C2S.OrderIDEx = proto.String(orderIDEx)
	}
	// send request and register receiving channel
	ch := make(chan *trdgetmaxtrdqtys.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetMaxTrdQtys, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	// wait for context signal or receiving result from channel then return to caller
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetMaxTrdQtys(), protocol.Error(resp)
	}
}
