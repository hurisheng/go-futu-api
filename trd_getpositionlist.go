package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdgetpositionlist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdGetPositionList = 2102 //Trd_GetPositionList	获取账户持仓
)

// 查询持仓
func (api *FutuAPI) GetPositionList(ctx context.Context, header *TrdHeader, filter *TrdFilterConditions, minPLRatio float64, maxPLRation float64, refresh bool) ([]*Position, error) {
	req := trdgetpositionlist.Request{
		C2S: &trdgetpositionlist.C2S{
			Header:           header.pb(),
			FilterConditions: filter.pb(),
			RefreshCache:     &refresh,
		},
	}
	if minPLRatio != 0 {
		req.C2S.FilterPLRatioMin = &minPLRatio
	}
	if maxPLRation != 0 {
		req.C2S.FilterPLRatioMax = &maxPLRation
	}
	ch := make(trdgetpositionlist.ResponseChan)
	if err := api.get(ProtoIDTrdGetPositionList, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return positionListFromPB(resp.GetS2C().GetPositionList()), protocol.Error(resp)
	}
}
