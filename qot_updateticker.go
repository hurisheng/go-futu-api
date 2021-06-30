package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdateticker"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateTicker = 3011 //Qot_UpdateTicker	推送逐笔
)

// 实时逐笔回调，异步处理已订阅股票的实时逐笔推送
func (api *FutuAPI) UpdateTicker(ctx context.Context) (<-chan *UpdateTickerResp, error) {
	ch := make(updateTickerChan)
	if err := api.update(ProtoIDQotUpdateTicker, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateTickerResp struct {
	Ticker *RTTicker
	Err    error
}

type updateTickerChan chan *UpdateTickerResp

var _ protocol.RespChan = make(updateTickerChan)

func (ch updateTickerChan) Send(b []byte) error {
	var resp qotupdateticker.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateTickerResp{
		Ticker: rtTickerFromUpdatePB(resp.GetS2C()),
		Err:    protocol.Error(&resp),
	}
	return nil
}
func (ch updateTickerChan) Close() {
	close(ch)
}

func rtTickerFromUpdatePB(pb *qotupdateticker.S2C) *RTTicker {
	if pb == nil {
		return nil
	}
	return &RTTicker{
		Security: securityFromPB(pb.GetSecurity()),
		Tickers:  tickerListFromPB(pb.GetTickerList()),
	}
}
