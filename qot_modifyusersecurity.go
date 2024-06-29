package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotmodifyusersecurity"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotModifyUserSecurity = 3214 //Qot_ModifyUserSecurity	修改自选股分组下的股票

func init() {
	workers[ProtoIDQotModifyUserSecurity] = protocol.NewGetter()
}

// 修改自选股列表
func (api *FutuAPI) ModifyUserSecurity(ctx context.Context, name string, op qotmodifyusersecurity.ModifyUserSecurityOp, securities []*qotcommon.Security) error {

	if name == "" || op == qotmodifyusersecurity.ModifyUserSecurityOp_ModifyUserSecurityOp_Unknown || len(securities) == 0 {
		return ErrParameters
	}
	req := &qotmodifyusersecurity.Request{
		C2S: &qotmodifyusersecurity.C2S{
			GroupName:    proto.String(name),
			Op:           proto.Int32(int32(op)),
			SecurityList: securities,
		},
	}
	ch := make(chan *qotmodifyusersecurity.Response)
	if err := api.proto.RegisterGet(ProtoIDQotModifyUserSecurity, req, protocol.NewProtobufChan(ch)); err != nil {
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
