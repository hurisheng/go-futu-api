package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetpricereminder"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetPriceReminder = 3221 //Qot_GetPriceReminder	获取到价提醒

func init() {
	workers[ProtoIDQotGetPriceReminder] = protocol.NewGetter()
}

// 获取到价提醒列表
func (api *FutuAPI) GetPriceReminder(ctx context.Context, security *qotcommon.Security, market qotcommon.QotMarket) ([]*qotgetpricereminder.PriceReminder, error) {

	if security == nil && market == qotcommon.QotMarket_QotMarket_Unknown {
		return nil, ErrParameters
	}
	req := &qotgetpricereminder.Request{
		C2S: &qotgetpricereminder.C2S{
			Security: security,
		},
	}
	if market != qotcommon.QotMarket_QotMarket_Unknown {
		req.C2S.Market = proto.Int32(int32(market))
	}
	ch := make(chan *qotgetpricereminder.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetPriceReminder, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetPriceReminderList(), protocol.Error(resp)
	}
}
