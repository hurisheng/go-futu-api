package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetreference"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetReference = 3206 //Qot_GetReference	获取正股相关股票
)

// 获取证券关联数据
func (api *FutuAPI) GetReferenceStockList(ctx context.Context, security *Security, refType qotgetreference.ReferenceType) ([]*SecurityStaticInfo, error) {
	req := qotgetreference.Request{
		C2S: &qotgetreference.C2S{
			Security:      security.pb(),
			ReferenceType: (*int32)(&refType),
		},
	}
	ch := make(qotgetreference.ResponseChan)
	if err := api.get(ProtoIDQotGetReference, &req, ch); err != nil {
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
