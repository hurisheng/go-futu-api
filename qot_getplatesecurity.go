package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetplatesecurity"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetPlateSecurity = 3205 //Qot_GetPlateSecurity	获取板块下的股票
)

// 获取板块内股票列表
func (api *FutuAPI) GetPlateStock(ctx context.Context, plate *Security, sortField qotcommon.SortField, ascend bool) ([]*SecurityStaticInfo, error) {
	req := qotgetplatesecurity.Request{
		C2S: &qotgetplatesecurity.C2S{
			Plate:  plate.pb(),
			Ascend: &ascend,
		},
	}
	if sortField != 0 {
		req.C2S.SortField = (*int32)(&sortField)
	}
	ch := make(qotgetplatesecurity.ResponseChan)
	if err := api.get(ProtoIDQotGetPlateSecurity, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return securityStaticInfoListFromPB(resp.GetS2C().GetStaticInfoList()), protocol.Error(resp)
	}
}
