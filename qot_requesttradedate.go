package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotrequesttradedate"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotRequestTradeDate = 3219 //Qot_RequestTradeDate	获取市场交易日，在线拉取不在本地计算
)

// 获取交易日
func (api *FutuAPI) RequestTradingDays(ctx context.Context, market qotcommon.TradeDateMarket, begin string, end string) ([]*TradeDate, error) {
	req := qotrequesttradedate.Request{
		C2S: &qotrequesttradedate.C2S{
			Market:    (*int32)(&market),
			BeginTime: &begin,
			EndTime:   &end,
		},
	}
	ch := make(qotrequesttradedate.ResponseChan)
	if err := api.get(ProtoIDQotRequestTradeDate, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return tradeDateListFromPB(resp.GetS2C().GetTradeDateList()), protocol.Error(resp)
	}
}

type TradeDate struct {
	Time          string                  //时间字符串
	Timestamp     float64                 //时间戳
	TradeDateType qotcommon.TradeDateType //Qot_Common.TradeDateType，交易时间类型
}

func tradeDateFromPB(pb *qotrequesttradedate.TradeDate) *TradeDate {
	if pb == nil {
		return nil
	}
	return &TradeDate{
		Time:          pb.GetTime(),
		Timestamp:     pb.GetTimestamp(),
		TradeDateType: qotcommon.TradeDateType(pb.GetTradeDateType()),
	}
}

func tradeDateListFromPB(pb []*qotrequesttradedate.TradeDate) []*TradeDate {
	if pb == nil {
		return nil
	}
	list := make([]*TradeDate, len(pb))
	for i, v := range pb {
		list[i] = tradeDateFromPB(v)
	}
	return list
}
