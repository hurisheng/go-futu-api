package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetownerplate"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetOwnerPlate = 3207 //Qot_GetOwnerPlate	获取股票所属板块
)

// 获取股票所属板块
func (api *FutuAPI) GetOwnerPlate(ctx context.Context, securities []*Security) ([]*OwnerPlate, error) {
	// 请求数据
	req := qotgetownerplate.Request{
		C2S: &qotgetownerplate.C2S{
			SecurityList: securityList(securities).pb(),
		},
	}
	ch := make(qotgetownerplate.ResponseChan)
	if err := api.get(ProtoIDQotGetOwnerPlate, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return ownerPlateListFromPB(resp.GetS2C().GetOwnerPlateList()), protocol.Error(resp)
	}
}

type OwnerPlate struct {
	Security   *Security
	PlateInfos []*PlateInfo
}

func ownerPlateFromPB(pb *qotgetownerplate.SecurityOwnerPlate) *OwnerPlate {
	if pb == nil {
		return nil
	}
	return &OwnerPlate{
		Security:   securityFromPB(pb.GetSecurity()),
		PlateInfos: plateInfoListFromPB(pb.GetPlateInfoList()),
	}
}

func ownerPlateListFromPB(pb []*qotgetownerplate.SecurityOwnerPlate) []*OwnerPlate {
	if pb == nil {
		return nil
	}
	p := make([]*OwnerPlate, len(pb))
	for i, v := range pb {
		p[i] = ownerPlateFromPB(v)
	}
	return p
}
