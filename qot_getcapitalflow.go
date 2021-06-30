package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetcapitalflow"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetCapitalFlow = 3211 //Qot_GetCapitalFlow	获取资金流向
)

//获取资金流向
func (api *FutuAPI) GetCapitalFlow(ctx context.Context, security *Security) (*CapitalFlow, error) {
	// 请求参数
	req := qotgetcapitalflow.Request{
		C2S: &qotgetcapitalflow.C2S{
			Security: security.pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetcapitalflow.ResponseChan)
	if err := api.get(ProtoIDQotGetCapitalFlow, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return capitalFlowFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type CapitalFlow struct {
	FlowItems          []*CapitalFlowItem //资金流向
	LastValidTIme      string             //数据最后有效时间字符串
	LastValidTimestamp float64            //数据最后有效时间戳
}

func capitalFlowFromPB(pb *qotgetcapitalflow.S2C) *CapitalFlow {
	if pb == nil {
		return nil
	}
	return &CapitalFlow{
		FlowItems:          flowItemListFromPB(pb.GetFlowItemList()),
		LastValidTIme:      pb.GetLastValidTime(),
		LastValidTimestamp: pb.GetLastValidTimestamp(),
	}
}

type CapitalFlowItem struct {
	InFlow    float64 //净流入的资金额度，正数代表流入，负数代表流出
	Time      string  //开始时间字符串，以分钟为单位
	Timestamp float64 //开始时间戳
}

func capitalFlowItemFromPB(pb *qotgetcapitalflow.CapitalFlowItem) *CapitalFlowItem {
	if pb == nil {
		return nil
	}
	return &CapitalFlowItem{
		InFlow:    pb.GetInFlow(),
		Time:      pb.GetTime(),
		Timestamp: pb.GetTimestamp(),
	}
}

func flowItemListFromPB(pb []*qotgetcapitalflow.CapitalFlowItem) []*CapitalFlowItem {
	if pb == nil {
		return nil
	}
	f := make([]*CapitalFlowItem, len(pb))
	for i, v := range pb {
		f[i] = capitalFlowItemFromPB(v)
	}
	return f
}
