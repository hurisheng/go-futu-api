package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdsubaccpush"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDTrdSubAccPush = 2008 //Trd_SubAccPush	订阅业务账户的交易推送数据

func init() {
	workers[ProtoIDTrdSubAccPush] = protocol.NewGetter()
}

// 订阅交易推送
func (api *FutuAPI) SubscribeTrd(ctx context.Context, accID []uint64) error {

	if len(accID) == 0 {
		return ErrParameters
	}
	req := &trdsubaccpush.Request{
		C2S: &trdsubaccpush.C2S{
			AccIDList: accID,
		},
	}

	ch := make(chan *trdsubaccpush.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdSubAccPush, req, protocol.NewProtobufChan(ch)); err != nil {
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
