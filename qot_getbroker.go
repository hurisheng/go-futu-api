package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetbroker"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetBroker = 3014 //Qot_GetBroker	获取经纪队列
)

// 获取实时经纪队列
func (api *FutuAPI) GetBrokerQueue(ctx context.Context, security *Security) (*BrokerQueue, error) {
	// 请求参数
	req := qotgetbroker.Request{
		C2S: &qotgetbroker.C2S{
			Security: security.pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetbroker.ResponseChan)
	if err := api.get(ProtoIDQotGetBroker, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return brokerQueueFromGetPB(resp.GetS2C()), protocol.Error(resp)
	}
}

func brokerQueueFromGetPB(pb *qotgetbroker.S2C) *BrokerQueue {
	if pb == nil {
		return nil
	}
	return &BrokerQueue{
		Security: securityFromPB(pb.GetSecurity()),
		Asks:     brokerListFromPB(pb.GetBrokerAskList()),
		Bids:     brokerListFromPB(pb.GetBrokerBidList()),
	}
}
