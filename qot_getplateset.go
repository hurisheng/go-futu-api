package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetplateset"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetPlateSet = 3204 //Qot_GetPlateSet	获取板块集合下的板块
)

// 获取板块列表
func (api *FutuAPI) GetPlateList(ctx context.Context, market qotcommon.QotMarket, plateClass qotcommon.PlateSetType) ([]*PlateInfo, error) {
	req := qotgetplateset.Request{
		C2S: &qotgetplateset.C2S{
			Market:       (*int32)(&market),
			PlateSetType: (*int32)(&plateClass),
		},
	}
	ch := make(qotgetplateset.ResponseChan)
	if err := api.get(ProtoIDQotGetPlateSet, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return plateInfoListFromPB(resp.GetS2C().GetPlateInfoList()), protocol.Error(resp)
	}
}
