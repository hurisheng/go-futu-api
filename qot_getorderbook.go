package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetorderbook"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetOrderBook = 3012 //Qot_GetOrderBook	获取买卖盘
)

// 获取实时摆盘
func (api *FutuAPI) GetOrderBook(ctx context.Context, security *Security, num int32) (*RTOrderBook, error) {
	// 请求参数
	req := qotgetorderbook.Request{
		C2S: &qotgetorderbook.C2S{
			Security: security.pb(),
			Num:      &num,
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetorderbook.ResponseChan)
	if err := api.get(ProtoIDQotGetOrderBook, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return rtOrderBookFromGetPB(resp.GetS2C()), protocol.Error(resp)
	}
}

func rtOrderBookFromGetPB(pb *qotgetorderbook.S2C) *RTOrderBook {
	if pb == nil {
		return nil
	}
	return &RTOrderBook{
		Security:                securityFromPB(pb.GetSecurity()),
		Asks:                    orderBookListFromPB(pb.GetOrderBookAskList()),
		Bids:                    orderBookListFromPB(pb.GetOrderBookBidList()),
		SvrRecvTimeBid:          pb.GetSvrRecvTimeBid(),
		SvrRecvTimeBidTimestamp: pb.GetSvrRecvTimeBidTimestamp(),
		SvrRecvTimeAsk:          pb.GetSvrRecvTimeAsk(),
		SvrRecvTimeAskTimestamp: pb.GetSvrRecvTimeAskTimestamp(),
	}
}
