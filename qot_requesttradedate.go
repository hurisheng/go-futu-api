package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotrequesttradedate"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotRequestTradeDate = 3219 //Qot_RequestTradeDate	获取市场交易日，在线拉取不在本地计算

func init() {
	workers[ProtoIDQotRequestTradeDate] = protocol.NewGetter()
}

// 获取交易日
func (api *FutuAPI) RequestTradingDays(ctx context.Context, market qotcommon.TradeDateMarket, begin string, end string,
	security *qotcommon.Security) ([]*qotrequesttradedate.TradeDate, error) {

	if begin == "" || end == "" {
		return nil, ErrParameters
	}
	req := &qotrequesttradedate.Request{
		C2S: &qotrequesttradedate.C2S{
			Market:    proto.Int32(int32(market)),
			BeginTime: proto.String(begin),
			EndTime:   proto.String(end),
			Security:  security,
		},
	}

	ch := make(chan *qotrequesttradedate.Response)
	if err := api.proto.RegisterGet(ProtoIDQotRequestTradeDate, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetTradeDateList(), protocol.Error(resp)
	}
}
