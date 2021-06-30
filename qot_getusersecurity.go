package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetusersecurity"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetUserSecurity = 3213 //Qot_GetUserSecurity	获取自选股分组下的股票
)

// 获取自选股列表
func (api *FutuAPI) GetUserSecurity(ctx context.Context, group string) ([]*SecurityStaticInfo, error) {
	req := qotgetusersecurity.Request{
		C2S: &qotgetusersecurity.C2S{
			GroupName: &group,
		},
	}
	ch := make(qotgetusersecurity.ResponseChan)
	if err := api.get(ProtoIDQotGetUserSecurity, &req, ch); err != nil {
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
