package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetticker"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetTicker = 3010 //Qot_GetTicker	获取逐笔
)

func (api *FutuAPI) GetRTTicker(ctx context.Context, security *Security, num int32) (*RTTicker, error) {
	// 请求参数
	req := qotgetticker.Request{C2S: &qotgetticker.C2S{
		Security:  security.pb(),
		MaxRetNum: &num,
	}}
	// 发送请求，同步返回结果
	ch := make(qotgetticker.ResponseChan)
	if err := api.get(ProtoIDQotGetTicker, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return rtTickerFromGetPB(resp.GetS2C()), protocol.Error(resp)
	}
}

func rtTickerFromGetPB(pb *qotgetticker.S2C) *RTTicker {
	if pb == nil {
		return nil
	}
	return &RTTicker{
		Security: securityFromPB(pb.GetSecurity()),
		Tickers:  tickerListFromPB(pb.GetTickerList()),
	}
}
