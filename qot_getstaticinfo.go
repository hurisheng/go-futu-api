package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetstaticinfo"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetStaticInfo = 3202 //Qot_GetStaticInfo	获取股票静态信息
)

// 获取静态数据
func (api *FutuAPI) GetStockBasicInfo(ctx context.Context, markte qotcommon.QotMarket, secType qotcommon.SecurityType, securities []*Security) ([]*SecurityStaticInfo, error) {
	req := qotgetstaticinfo.Request{
		C2S: &qotgetstaticinfo.C2S{
			Market:       (*int32)(&markte),
			SecType:      (*int32)(&secType),
			SecurityList: securityList(securities).pb(),
		},
	}
	ch := make(qotgetstaticinfo.ResponseChan)
	if err := api.get(ProtoIDQotGetStaticInfo, &req, ch); err != nil {
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
