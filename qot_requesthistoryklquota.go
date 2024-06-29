package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotrequesthistoryklquota"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotRequestHistoryKLQuota = 3104 //Qot_RequestHistoryKLQuota	获取历史 K 线额度

func init() {
	workers[ProtoIDQotRequestHistoryKLQuota] = protocol.NewGetter()
}

func (api *FutuAPI) GetHistoryKLQuota(ctx context.Context, detail bool) (*qotrequesthistoryklquota.S2C, error) {
	req := &qotrequesthistoryklquota.Request{
		C2S: &qotrequesthistoryklquota.C2S{
			BGetDetail: proto.Bool(detail),
		},
	}
	ch := make(chan *qotrequesthistoryklquota.Response)
	if err := api.proto.RegisterGet(ProtoIDQotRequestHistoryKLQuota, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), protocol.Error(resp)
	}
}
