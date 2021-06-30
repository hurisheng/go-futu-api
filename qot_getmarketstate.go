package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetmarketstate"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetMarketState = 3223 //Qot_GetMarketState	获取指定品种的市场状态
)

//获取标的市场状态
func (api *FutuAPI) GetMarketState(ctx context.Context, securities []*Security) ([]*MarketInfo, error) {
	// 请求参数
	req := qotgetmarketstate.Request{
		C2S: &qotgetmarketstate.C2S{
			SecurityList: securityList(securities).pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetmarketstate.ResponseChan)
	if err := api.get(ProtoIDQotGetMarketState, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return marketInfoListFromPB(resp.GetS2C().GetMarketInfoList()), protocol.Error(resp)
	}
}

func marketInfoListFromPB(pb []*qotgetmarketstate.MarketInfo) []*MarketInfo {
	if pb == nil {
		return nil
	}
	m := make([]*MarketInfo, len(pb))
	for i, v := range pb {
		m[i] = marketInfoFromPB(v)
	}
	return m
}

type MarketInfo struct {
	Security    *Security                //股票代码
	Name        string                   // 股票名称
	MarketState qotcommon.QotMarketState //Qot_Common.QotMarketState，市场状态
}

func marketInfoFromPB(pb *qotgetmarketstate.MarketInfo) *MarketInfo {
	if pb == nil {
		return nil
	}
	return &MarketInfo{
		Security:    securityFromPB(pb.GetSecurity()),
		Name:        pb.GetName(),
		MarketState: qotcommon.QotMarketState(pb.GetMarketState()),
	}
}
