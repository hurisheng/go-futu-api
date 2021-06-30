package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetoptionchain"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetOptionChain = 3209 //Qot_GetOptionChain	获取期权链
)

// 获取期权链
func (api *FutuAPI) GetOptionChain(ctx context.Context, owner *Security, begin string, end string,
	indexOptType qotcommon.IndexOptionType, optType qotcommon.OptionType, cond qotgetoptionchain.OptionCondType, filter *DataFilter) ([]*OptionChain, error) {
	// 请求参数
	req := qotgetoptionchain.Request{
		C2S: &qotgetoptionchain.C2S{
			Owner:           owner.pb(),
			IndexOptionType: (*int32)(&indexOptType),
			Type:            (*int32)(&optType),
			Condition:       (*int32)(&cond),
			BeginTime:       &begin,
			EndTime:         &end,
			DataFilter:      filter.pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetoptionchain.ResponseChan)
	if err := api.get(ProtoIDQotGetOptionChain, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return optionChainListFromPB(resp.GetS2C().OptionChain), protocol.Error(resp)
	}
}

type OptionChain struct {
	StrikeTime      string        //行权日
	Option          []*OptionItem //期权信息
	StrikeTimestamp float64       //行权日时间戳
}

func optionChainFromPB(pb *qotgetoptionchain.OptionChain) *OptionChain {
	if pb == nil {
		return nil
	}
	return &OptionChain{
		StrikeTime:      pb.GetStrikeTime(),
		Option:          optionItemListFromPB(pb.GetOption()),
		StrikeTimestamp: pb.GetStrikeTimestamp(),
	}
}

func optionChainListFromPB(pb []*qotgetoptionchain.OptionChain) []*OptionChain {
	if pb == nil {
		return nil
	}
	list := make([]*OptionChain, len(pb))
	for i, v := range pb {
		list[i] = optionChainFromPB(v)
	}
	return list
}

type OptionItem struct {
	Call *SecurityStaticInfo //看涨期权，不一定有该字段，由请求条件决定
	Put  *SecurityStaticInfo //看跌期权，不一定有该字段，由请求条件决定
}

func optionItemFromPB(pb *qotgetoptionchain.OptionItem) *OptionItem {
	if pb == nil {
		return nil
	}
	return &OptionItem{
		Call: securityStaticInfoFromPB(pb.GetCall()),
		Put:  securityStaticInfoFromPB(pb.GetPut()),
	}
}

func optionItemListFromPB(pb []*qotgetoptionchain.OptionItem) []*OptionItem {
	if pb == nil {
		return nil
	}
	list := make([]*OptionItem, len(pb))
	for i, v := range pb {
		list[i] = optionItemFromPB(v)
	}
	return list
}

type FilterDouble struct {
	Value float64
}

func (d *FilterDouble) pb() *float64 {
	if d == nil {
		return nil
	}
	return &d.Value
}

type DataFilter struct {
	ImpliedVolatilityMin *FilterDouble //隐含波动率过滤
	ImpliedVolatilityMax *FilterDouble
	DeltaMin             *FilterDouble
	DeltaMax             *FilterDouble
	GammaMin             *FilterDouble
	GammaMax             *FilterDouble
	VegaMin              *FilterDouble
	VegaMax              *FilterDouble
	ThetaMin             *FilterDouble
	ThetaMax             *FilterDouble
	RhoMin               *FilterDouble
	RhoMax               *FilterDouble
	NetOpenInterestMin   *FilterDouble //净未平仓合约数过滤
	NetOpenInterestMax   *FilterDouble
	OpenInterestMin      *FilterDouble //未平仓合约数过滤
	OpenInterestMax      *FilterDouble
	VolMin               *FilterDouble //成交量过滤
	VolMax               *FilterDouble
}

func (f *DataFilter) pb() *qotgetoptionchain.DataFilter {
	if f == nil {
		return nil
	}
	var d qotgetoptionchain.DataFilter
	if f.ImpliedVolatilityMin != nil {
		d.ImpliedVolatilityMin = &f.ImpliedVolatilityMin.Value
	}
	if f.ImpliedVolatilityMax != nil {
		d.ImpliedVolatilityMax = &f.ImpliedVolatilityMax.Value
	}
	if f.DeltaMin != nil {
		d.DeltaMin = &f.DeltaMin.Value
	}
	if f.DeltaMax != nil {
		d.DeltaMax = &f.DeltaMax.Value
	}
	if f.GammaMin != nil {
		d.GammaMin = &f.DeltaMin.Value
	}
	if f.GammaMax != nil {
		d.GammaMax = &f.DeltaMax.Value
	}
	if f.VegaMin != nil {
		d.VegaMin = &f.VegaMin.Value
	}
	if f.VegaMax != nil {
		d.VegaMax = &f.VegaMax.Value
	}
	if f.ThetaMin != nil {
		d.ThetaMin = &f.ThetaMin.Value
	}
	if f.ThetaMax != nil {
		d.ThetaMax = &f.ThetaMax.Value
	}
	if f.RhoMin != nil {
		d.RhoMin = &f.RhoMin.Value
	}
	if f.RhoMax != nil {
		d.RhoMax = &f.RhoMax.Value
	}
	if f.NetOpenInterestMin != nil {
		d.NetOpenInterestMin = &f.NetOpenInterestMin.Value
	}
	if f.NetOpenInterestMax != nil {
		d.NetOpenInterestMax = &f.NetOpenInterestMax.Value
	}
	if f.OpenInterestMin != nil {
		d.OpenInterestMin = &f.OpenInterestMin.Value
	}
	if f.OpenInterestMax != nil {
		d.OpenInterestMax = &f.OpenInterestMax.Value
	}
	if f.VolMin != nil {
		d.VolMin = &f.VolMin.Value
	}
	if f.VolMax != nil {
		d.VolMax = &f.VolMax.Value
	}
	return &d
}
