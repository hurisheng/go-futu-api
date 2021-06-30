package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetsubinfo"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetSubInfo = 3003 //Qot_GetSubInfo	获取订阅信息
)

// 获取订阅信息
func (api *FutuAPI) QuerySubscription(ctx context.Context, isAll bool) (*Subscription, error) {
	// 请求参数
	req := qotgetsubinfo.Request{C2S: &qotgetsubinfo.C2S{
		IsReqAllConn: &isAll,
	}}
	// 发送请求，同步返回结果
	ch := make(qotgetsubinfo.ResponseChan)
	if err := api.get(ProtoIDQotGetSubInfo, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return subscriptionFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type Subscription struct {
	ConnSubInfos   []*ConnSubInfo //单条连接订阅信息
	TotalUsedQuota int32          //*FutuOpenD 已使用的订阅额度
	RemainQuota    int32          //*FutuOpenD 剩余订阅额度
}

func subscriptionFromPB(pb *qotgetsubinfo.S2C) *Subscription {
	if pb == nil {
		return nil
	}
	return &Subscription{
		TotalUsedQuota: pb.GetTotalUsedQuota(),
		RemainQuota:    pb.GetRemainQuota(),
		ConnSubInfos:   connSubInfoListFromPB(pb.GetConnSubInfoList()),
	}
}
