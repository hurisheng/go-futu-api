package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdateorderbook"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateOrderBook = 3013 //Qot_UpdateOrderBook	推送买卖盘
)

// 实时摆盘回调
func (api *FutuAPI) UpdateOrderBook(ctx context.Context) (<-chan *UpdateOrderBookResp, error) {
	ch := make(updateOrderBookChan)
	if err := api.update(ProtoIDQotUpdateOrderBook, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateOrderBookResp struct {
	OrderBook *RTOrderBook
	Err       error
}

type updateOrderBookChan chan *UpdateOrderBookResp

var _ protocol.RespChan = make(updateOrderBookChan)

func (ch updateOrderBookChan) Send(b []byte) error {
	var resp qotupdateorderbook.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateOrderBookResp{
		OrderBook: rtOrderBookFromUpdatePB(resp.GetS2C()),
		Err:       protocol.Error(&resp),
	}
	return nil
}

func (ch updateOrderBookChan) Close() {
	close(ch)
}

func rtOrderBookFromUpdatePB(pb *qotupdateorderbook.S2C) *RTOrderBook {
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
