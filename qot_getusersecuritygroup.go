package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetusersecuritygroup"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetUserSecurityGroup = 3222 //Qot_GetUserSecurityGroup	获取自选股分组列表

func init() {
	workers[ProtoIDQotGetUserSecurityGroup] = protocol.NewGetter()
}

// 获取自选股分组
func (api *FutuAPI) GetUserSecurityGroup(ctx context.Context, groupType qotgetusersecuritygroup.GroupType) ([]*qotgetusersecuritygroup.GroupData, error) {

	if groupType == qotgetusersecuritygroup.GroupType_GroupType_Unknown {
		return nil, ErrParameters
	}
	req := &qotgetusersecuritygroup.Request{
		C2S: &qotgetusersecuritygroup.C2S{
			GroupType: proto.Int32(int32(groupType)),
		},
	}
	ch := make(chan *qotgetusersecuritygroup.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetUserSecurityGroup, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetGroupList(), protocol.Error(resp)
	}
}
