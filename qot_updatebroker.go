package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatebroker"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateBroker = 3015 //Qot_UpdateBroker	推送经纪队列
)

// 实时经纪队列回调
func (api *FutuAPI) UpdateBroker(ctx context.Context) (<-chan *UpdateBrokerResp, error) {
	ch := make(updateBrokerChan)
	if err := api.update(ProtoIDQotUpdateBroker, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateBrokerResp struct {
	BrokerQueue *BrokerQueue
	Err         error
}

type updateBrokerChan chan *UpdateBrokerResp

var _ protocol.RespChan = make(updateBrokerChan)

func (ch updateBrokerChan) Send(b []byte) error {
	var resp qotupdatebroker.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateBrokerResp{
		BrokerQueue: brokerQueueFromUpdatePB(resp.GetS2C()),
		Err:         protocol.Error(&resp),
	}
	return nil
}

func (ch updateBrokerChan) Close() {
	close(ch)
}

func brokerQueueFromUpdatePB(pb *qotupdatebroker.S2C) *BrokerQueue {
	if pb == nil {
		return nil
	}
	return &BrokerQueue{
		Security: securityFromPB(pb.GetSecurity()),
		Asks:     brokerListFromPB(pb.GetBrokerAskList()),
		Bids:     brokerListFromPB(pb.GetBrokerBidList()),
	}
}
