package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetusersecuritygroup"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetUserSecurityGroup = 3222 //Qot_GetUserSecurityGroup	获取自选股分组列表
)

// 获取自选股分组
func (api *FutuAPI) GetUserSecurityGroup(ctx context.Context, groupType qotgetusersecuritygroup.GroupType) ([]*GroupData, error) {
	req := qotgetusersecuritygroup.Request{
		C2S: &qotgetusersecuritygroup.C2S{
			GroupType: (*int32)(&groupType),
		},
	}
	ch := make(qotgetusersecuritygroup.ResponseChan)
	if err := api.get(ProtoIDQotGetUserSecurityGroup, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return groupDataListFromPB(resp.GetS2C().GetGroupList()), protocol.Error(resp)
	}
}

type GroupData struct {
	Name string                            // 自选股分组名字
	Type qotgetusersecuritygroup.GroupType // GroupType，自选股分组类型。
}

func groupDataFromPB(pb *qotgetusersecuritygroup.GroupData) *GroupData {
	if pb == nil {
		return nil
	}
	return &GroupData{
		Name: pb.GetGroupName(),
		Type: qotgetusersecuritygroup.GroupType(pb.GetGroupType()),
	}
}

func groupDataListFromPB(pb []*qotgetusersecuritygroup.GroupData) []*GroupData {
	if pb == nil {
		return nil
	}
	list := make([]*GroupData, len(pb))
	for i, v := range pb {
		list[i] = groupDataFromPB(v)
	}
	return list
}
