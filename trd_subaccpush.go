package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdsubaccpush"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdSubAccPush = 2008 //Trd_SubAccPush	订阅业务账户的交易推送数据
)

// 订阅交易推送
func (api *FutuAPI) SubscribeTrd(ctx context.Context, accID []uint64) error {
	req := trdsubaccpush.Request{
		C2S: &trdsubaccpush.C2S{
			AccIDList: accID,
		},
	}
	ch := make(trdsubaccpush.ResponseChan)
	if err := api.get(ProtoIDTrdSubAccPush, &req, ch); err != nil {
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
