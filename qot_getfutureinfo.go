package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetfutureinfo"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetFutureInfo = 3218 //Qot_GetFutureInfo	获取期货合约资料
)

// 获取期货合约资料
func (api *FutuAPI) GetFutureInfo(ctx context.Context, securities []*Security) ([]*FutureInfo, error) {
	req := qotgetfutureinfo.Request{
		C2S: &qotgetfutureinfo.C2S{
			SecurityList: securityList(securities).pb(),
		},
	}
	ch := make(qotgetfutureinfo.ResponseChan)
	if err := api.get(ProtoIDQotGetFutureInfo, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return futureInfoListFromPB(resp.GetS2C().GetFutureInfoList()), protocol.Error(resp)
	}
}

type FutureInfo struct {
	Name               string       // 合约名称
	Security           *Security    // 合约代码
	LastTradeTime      string       //最后交易日，只有非主连期货合约才有该字段
	LastTradeTimestamp float64      //最后交易日时间戳，只有非主连期货合约才有该字段
	Owner              *Security    //标的股 股票期货和股指期货才有该字段
	OwnerOther         string       //标的
	Exchange           string       //交易所
	ContractType       string       //合约类型
	ContractSize       float64      //合约规模
	ContractSizeUnit   string       //合约规模的单位
	QuoteCurrency      string       //报价货币
	MinVar             float64      //最小变动单位
	MinVarUnit         string       //最小变动单位的单位
	QuoteUnit          string       //报价单位
	TradeTime          []*TradeTime //交易时间
	TimeZone           string       //所在时区
	ExchangeFormatURL  string       //交易所规格
}

func futureInfoFromPB(pb *qotgetfutureinfo.FutureInfo) *FutureInfo {
	if pb == nil {
		return nil
	}
	return &FutureInfo{
		Name:               pb.GetName(),
		Security:           securityFromPB(pb.GetSecurity()),
		LastTradeTime:      pb.GetLastTradeTime(),
		LastTradeTimestamp: pb.GetLastTradeTimestamp(),
		Owner:              securityFromPB(pb.GetOwner()),
		OwnerOther:         pb.GetOwnerOther(),
		Exchange:           pb.GetExchange(),
		ContractType:       pb.GetContractType(),
		ContractSize:       pb.GetContractSize(),
		ContractSizeUnit:   pb.GetContractSizeUnit(),
		QuoteCurrency:      pb.GetQuoteCurrency(),
		MinVar:             pb.GetMinVar(),
		MinVarUnit:         pb.GetMinVarUnit(),
		QuoteUnit:          pb.GetQuoteUnit(),
		TradeTime:          tradeTimeListFromPB(pb.GetTradeTime()),
		TimeZone:           pb.GetTimeZone(),
		ExchangeFormatURL:  pb.GetExchangeFormatUrl(),
	}
}

func futureInfoListFromPB(pb []*qotgetfutureinfo.FutureInfo) []*FutureInfo {
	if pb == nil {
		return nil
	}
	list := make([]*FutureInfo, len(pb))
	for i, v := range pb {
		list[i] = futureInfoFromPB(v)
	}
	return list
}

type TradeTime struct {
	Begin float64 // 开始时间，以分钟为单位
	End   float64 // 结束时间，以分钟为单位
}

func tradeTimeFromPB(pb *qotgetfutureinfo.TradeTime) *TradeTime {
	if pb == nil {
		return nil
	}
	return &TradeTime{
		Begin: pb.GetBegin(),
		End:   pb.GetEnd(),
	}
}

func tradeTimeListFromPB(pb []*qotgetfutureinfo.TradeTime) []*TradeTime {
	if pb == nil {
		return nil
	}
	list := make([]*TradeTime, len(pb))
	for i, v := range pb {
		list[i] = tradeTimeFromPB(v)
	}
	return list
}
