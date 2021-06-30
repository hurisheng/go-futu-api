package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotmodifyusersecurity"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotModifyUserSecurity = 3214 //Qot_ModifyUserSecurity	修改自选股分组下的股票
)

// 修改自选股列表
func (api *FutuAPI) ModifyUserSecurity(ctx context.Context, name string, op qotmodifyusersecurity.ModifyUserSecurityOp, securities []*Security) error {
	req := qotmodifyusersecurity.Request{
		C2S: &qotmodifyusersecurity.C2S{
			GroupName:    &name,
			Op:           (*int32)(&op),
			SecurityList: securityList(securities).pb(),
		},
	}
	ch := make(qotmodifyusersecurity.ResponseChan)
	if err := api.get(ProtoIDQotModifyUserSecurity, &req, ch); err != nil {
		return nil
	}
	select {
	case <-ctx.Done():
		return ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return ErrChannelClosed
		}
		return protocol.Error(resp)
	}
}
