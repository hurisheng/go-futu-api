package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetmaxtrdqtys"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdGetMaxTrdQtys = 2111 //Trd_GetMaxTrdQtys	获取最大交易数量
)

// 查询最大可买可卖
func (api *FutuAPI) GetMaxTrdQtys(ctx context.Context, header *TrdHeader, orderType trdcommon.OrderType, code string, price float64,
	orderID uint64, adjust bool, sideAndLimit float64, secMarket trdcommon.TrdSecMarket) (*MaxTrdQtys, error) {
	req := trdgetmaxtrdqtys.Request{
		C2S: &trdgetmaxtrdqtys.C2S{
			Header:             header.pb(),
			OrderType:          (*int32)(&orderType),
			Code:               &code,
			Price:              &price,
			AdjustPrice:        &adjust,
			AdjustSideAndLimit: &sideAndLimit,
		},
	}
	if orderID != 0 {
		req.C2S.OrderID = &orderID
	}
	if secMarket != 0 {
		req.C2S.SecMarket = (*int32)(&secMarket)
	}
	ch := make(trdgetmaxtrdqtys.ResponseChan)
	if err := api.get(ProtoIDTrdGetMaxTrdQtys, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return maxTrdQtysFromPB(resp.GetS2C().GetMaxTrdQtys()), protocol.Error(resp)
	}
}
