package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetusersecurity"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetUserSecurity = 3213 //Qot_GetUserSecurity	获取自选股分组下的股票

func init() {
	workers[ProtoIDQotGetUserSecurity] = protocol.NewGetter()
}

// 获取自选股列表
func (api *FutuAPI) GetUserSecurity(ctx context.Context, group string) ([]*qotcommon.SecurityStaticInfo, error) {

	if group == "" {
		return nil, ErrParameters
	}
	req := &qotgetusersecurity.Request{
		C2S: &qotgetusersecurity.C2S{
			GroupName: proto.String(group),
		},
	}
	ch := make(chan *qotgetusersecurity.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetUserSecurity, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetStaticInfoList(), protocol.Error(resp)
	}
}
