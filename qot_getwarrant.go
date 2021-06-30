package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetwarrant"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetWarrant = 3210 //Qot_GetWarrant	获取窝轮
)

// 筛选窝轮
func (api *FutuAPI) GetWarrant(ctx context.Context, begin int32, num int32, sortField qotcommon.SortField, ascend bool, filter *WarrantFilter) (*Warrant, error) {
	req := qotgetwarrant.Request{
		C2S: &qotgetwarrant.C2S{
			Begin:     &begin,
			Num:       &num,
			SortField: (*int32)(&sortField),
			Ascend:    &ascend,
		},
	}
	filter.pb(&req)
	ch := make(qotgetwarrant.ResponseChan)
	if err := api.get(ProtoIDQotGetWarrant, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return warrantFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type FilterString struct {
	Value string
}

type FilterInt32 struct {
	Value int32
}

type FilterUInt64 struct {
	Value uint64
}

type WarrantFilter struct {
	Owner                 *Security               //所属正股
	TypeList              []qotcommon.WarrantType //Qot_Common.WarrantType，窝轮类型过滤列表
	IssuerList            []qotcommon.Issuer      //Qot_Common.Issuer，发行人过滤列表
	IpoPeriod             qotcommon.IpoPeriod     //Qot_Common.IpoPeriod，上市日
	PriceType             qotcommon.PriceType     //Qot_Common.PriceType，价内/价外（暂不支持界内证的界内外筛选）
	Status                qotcommon.WarrantStatus //Qot_Common.WarrantStatus，窝轮状态
	MaturityTimeMin       *FilterString           //到期日，到期日范围的开始时间戳
	MaturityTimeMax       *FilterString           //到期日范围的结束时间戳
	CurPriceMin           *FilterDouble           //最新价的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	CurPriceMax           *FilterDouble           //最新价的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	StrikePriceMin        *FilterDouble           //行使价的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	StrikePriceMax        *FilterDouble           //行使价的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	StreetMin             *FilterDouble           //街货占比的过滤下限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	StreetMax             *FilterDouble           //街货占比的过滤上限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	ConversionMin         *FilterDouble           //换股比率的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	ConversionMax         *FilterDouble           //换股比率的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	VolMin                *FilterUInt64           //成交量的过滤下限（闭区间），不传代表下限为 -∞
	VolMax                *FilterUInt64           //成交量的过滤上限（闭区间），不传代表上限为 +∞
	PremiumMin            *FilterDouble           //溢价的过滤下限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	PremiumMax            *FilterDouble           //溢价的过滤上限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	LeverageRatioMin      *FilterDouble           //杠杆比率的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	LeverageRatioMax      *FilterDouble           //杠杆比率的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	DeltaMin              *FilterDouble           //对冲值的过滤下限（闭区间），仅认购认沽支持此字段过滤，不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	DeltaMax              *FilterDouble           //对冲值的过滤上限（闭区间），仅认购认沽支持此字段过滤，不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	ImplieMin             *FilterDouble           //引伸波幅的过滤下限（闭区间），仅认购认沽支持此字段过滤，不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	ImplieMax             *FilterDouble           //引伸波幅的过滤上限（闭区间），仅认购认沽支持此字段过滤，不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	RecoveryPriceMin      *FilterDouble           //收回价的过滤下限（闭区间），仅牛熊证支持此字段过滤，不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	RecoveryPriceMax      *FilterDouble           //收回价的过滤上限（闭区间），仅牛熊证支持此字段过滤，不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	PriceRecoveryRatioMin *FilterDouble           //正股距收回价，的过滤下限（闭区间），仅牛熊证支持此字段过滤。该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	PriceRecoveryRatioMax *FilterDouble           //正股距收回价，的过滤上限（闭区间），仅牛熊证支持此字段过滤。该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
}

func (filter *WarrantFilter) pb(req *qotgetwarrant.Request) {
	if filter != nil {
		if filter.Owner != nil {
			req.C2S.Owner = filter.Owner.pb()
		}
		if filter.TypeList != nil {
			req.C2S.TypeList = make([]int32, len(filter.TypeList))
			for i, v := range filter.TypeList {
				req.C2S.TypeList[i] = int32(v)
			}
		}
		if filter.IssuerList != nil {
			req.C2S.IssuerList = make([]int32, len(filter.IssuerList))
			for i, v := range filter.IssuerList {
				req.C2S.IssuerList[i] = int32(v)
			}
		}
		if filter.IpoPeriod != 0 {
			req.C2S.IpoPeriod = (*int32)(&filter.IpoPeriod)
		}
		if filter.PriceType != 0 {
			req.C2S.PriceType = (*int32)(&filter.PriceType)
		}
		if filter.Status != 0 {
			req.C2S.Status = (*int32)(&filter.Status)
		}
		if filter.MaturityTimeMin != nil {
			req.C2S.MaturityTimeMin = &filter.MaturityTimeMin.Value
		}
		if filter.MaturityTimeMax != nil {
			req.C2S.MaturityTimeMax = &filter.MaturityTimeMax.Value
		}
		if filter.CurPriceMin != nil {
			req.C2S.CurPriceMin = &filter.CurPriceMin.Value
		}
		if filter.CurPriceMax != nil {
			req.C2S.CurPriceMax = &filter.CurPriceMax.Value
		}
		if filter.StrikePriceMin != nil {
			req.C2S.StrikePriceMin = &filter.StrikePriceMin.Value
		}
		if filter.StrikePriceMax != nil {
			req.C2S.StrikePriceMax = &filter.StrikePriceMax.Value
		}
		if filter.StreetMin != nil {
			req.C2S.StreetMin = &filter.StreetMin.Value
		}
		if filter.StreetMax != nil {
			req.C2S.StreetMax = &filter.StreetMax.Value
		}
		if filter.ConversionMin != nil {
			req.C2S.ConversionMin = &filter.ConversionMin.Value
		}
		if filter.ConversionMax != nil {
			req.C2S.ConversionMax = &filter.ConversionMax.Value
		}
		if filter.VolMin != nil {
			req.C2S.VolMin = &filter.VolMin.Value
		}
		if filter.VolMax != nil {
			req.C2S.VolMax = &filter.VolMax.Value
		}
		if filter.VolMin != nil {
			req.C2S.VolMin = &filter.VolMin.Value
		}
		if filter.PremiumMin != nil {
			req.C2S.PremiumMin = &filter.PremiumMin.Value
		}
		if filter.PremiumMax != nil {
			req.C2S.PremiumMax = &filter.PremiumMax.Value
		}
		if filter.LeverageRatioMin != nil {
			req.C2S.LeverageRatioMin = &filter.LeverageRatioMin.Value
		}
		if filter.LeverageRatioMax != nil {
			req.C2S.LeverageRatioMax = &filter.LeverageRatioMax.Value
		}
		if filter.DeltaMin != nil {
			req.C2S.DeltaMin = &filter.DeltaMin.Value
		}
		if filter.DeltaMax != nil {
			req.C2S.DeltaMax = &filter.DeltaMax.Value
		}
		if filter.ImplieMin != nil {
			req.C2S.ImpliedMin = &filter.ImplieMin.Value
		}
		if filter.ImplieMax != nil {
			req.C2S.ImpliedMax = &filter.ImplieMax.Value
		}
		if filter.RecoveryPriceMin != nil {
			req.C2S.RecoveryPriceMin = &filter.RecoveryPriceMin.Value
		}
		if filter.RecoveryPriceMax != nil {
			req.C2S.RecoveryPriceMax = &filter.RecoveryPriceMax.Value
		}
		if filter.PriceRecoveryRatioMin != nil {
			req.C2S.PriceRecoveryRatioMin = &filter.PriceRecoveryRatioMin.Value
		}
		if filter.PriceRecoveryRatioMax != nil {
			req.C2S.PriceRecoveryRatioMax = &filter.PriceRecoveryRatioMax.Value
		}
	}
}

type Warrant struct {
	LastPage    bool           //是否最后一页了，false:非最后一页，还有窝轮记录未返回; true:已是最后一页
	AllCount    int32          //该条件请求所有数据的个数
	WarrantList []*WarrantData //窝轮数据
}

func warrantFromPB(pb *qotgetwarrant.S2C) *Warrant {
	if pb == nil {
		return nil
	}
	return &Warrant{
		LastPage:    pb.GetLastPage(),
		AllCount:    pb.GetAllCount(),
		WarrantList: warrantDataListFromPB(pb.GetWarrantDataList()),
	}
}

type WarrantData struct {
	// 静态数据项
	Stock              *Security             //股票
	Owner              *Security             //所属正股
	Type               qotcommon.WarrantType //Qot_Common.WarrantType，窝轮类型
	Issuer             qotcommon.Issuer      //Qot_Common.Issuer，发行人
	MaturityTime       string                //到期日
	MaturityTimestamp  float64               //到期日时间戳
	ListTime           string                //上市时间
	ListTimestamp      float64               //上市时间戳
	LastTradeTime      string                //最后交易日
	LastTradeTimestamp float64               //最后交易日时间戳
	RecoveryPrice      float64               //收回价，仅牛熊证支持此字段
	ConversionRatio    float64               //换股比率
	LotSize            int32                 //每手数量
	StrikePrice        float64               //行使价
	LastClosePrice     float64               //昨收价
	Name               string                //名称

	// 动态数据项
	CurPrice           float64                 //当前价
	PriceChangeVal     float64                 //涨跌额
	ChangeRate         float64                 //涨跌幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	Status             qotcommon.WarrantStatus //Qot_Common.WarrantStatus，窝轮状态
	BidPrice           float64                 //买入价
	AskPrice           float64                 //卖出价
	BidVol             int64                   //买量
	AskVol             int64                   //卖量
	Volume             int64                   //成交量
	Turnover           float64                 //成交额
	Score              float64                 //综合评分
	Premium            float64                 //溢价（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	BreakEvenPoint     float64                 //打和点
	Leverage           float64                 //杠杆比率（倍）
	IPOP               float64                 //价内/价外，正数表示价内，负数表示价外（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	PriceRecoveryRatio float64                 //正股距收回价，仅牛熊证支持此字段（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	ConversionPrice    float64                 //换股价
	StreetRate         float64                 //街货占比（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	StreetVol          int64                   //街货量
	Amplitude          float64                 //振幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	IssueSize          int64                   //发行量
	HighPrice          float64                 //最高价
	LowPrice           float64                 //最低价
	ImpliedVolatility  float64                 //引申波幅，仅认购认沽支持此字段
	Delta              float64                 //对冲值，仅认购认沽支持此字段
	EffectiveLeverage  float64                 //有效杠杆
	UpperStrikePrice   float64                 //上限价，仅界内证支持此字段
	LowerStrikePrice   float64                 //下限价，仅界内证支持此字段
	InLinePriceStatus  qotcommon.PriceType     //Qot_Common.PriceType，界内界外，仅界内证支持此字段
}

func warrantDataFromPB(pb *qotgetwarrant.WarrantData) *WarrantData {
	if pb == nil {
		return nil
	}
	return &WarrantData{
		Stock:              securityFromPB(pb.GetStock()),
		Owner:              securityFromPB(pb.GetOwner()),
		Type:               qotcommon.WarrantType(pb.GetType()),
		Issuer:             qotcommon.Issuer(pb.GetIssuer()),
		MaturityTime:       pb.GetMaturityTime(),
		MaturityTimestamp:  pb.GetMaturityTimestamp(),
		ListTime:           pb.GetListTime(),
		ListTimestamp:      pb.GetListTimestamp(),
		LastTradeTime:      pb.GetLastTradeTime(),
		LastTradeTimestamp: pb.GetLastTradeTimestamp(),
		RecoveryPrice:      pb.GetRecoveryPrice(),
		ConversionRatio:    pb.GetConversionRatio(),
		LotSize:            pb.GetLotSize(),
		StrikePrice:        pb.GetStrikePrice(),
		LastClosePrice:     pb.GetLastClosePrice(),
		Name:               pb.GetName(),
		CurPrice:           pb.GetCurPrice(),
		PriceChangeVal:     pb.GetPriceChangeVal(),
		ChangeRate:         pb.GetChangeRate(),
		Status:             qotcommon.WarrantStatus(pb.GetStatus()),
		BidPrice:           pb.GetBidPrice(),
		AskPrice:           pb.GetAskPrice(),
		BidVol:             pb.GetBidVol(),
		AskVol:             pb.GetAskVol(),
		Volume:             pb.GetVolume(),
		Turnover:           pb.GetTurnover(),
		Score:              pb.GetScore(),
		Premium:            pb.GetPremium(),
		BreakEvenPoint:     pb.GetBreakEvenPoint(),
		Leverage:           pb.GetLeverage(),
		IPOP:               pb.GetIpop(),
		PriceRecoveryRatio: pb.GetPriceRecoveryRatio(),
		ConversionPrice:    pb.GetConversionPrice(),
		StreetRate:         pb.GetStreetRate(),
		StreetVol:          pb.GetStreetVol(),
		Amplitude:          pb.GetAmplitude(),
		IssueSize:          pb.GetIssueSize(),
		HighPrice:          pb.GetHighPrice(),
		LowPrice:           pb.GetLowPrice(),
		ImpliedVolatility:  pb.GetImpliedVolatility(),
		Delta:              pb.GetDelta(),
		EffectiveLeverage:  pb.GetEffectiveLeverage(),
		UpperStrikePrice:   pb.GetUpperStrikePrice(),
		LowerStrikePrice:   pb.GetLowerStrikePrice(),
		InLinePriceStatus:  qotcommon.PriceType(pb.GetInLinePriceStatus()),
	}
}

func warrantDataListFromPB(pb []*qotgetwarrant.WarrantData) []*WarrantData {
	if pb == nil {
		return nil
	}
	list := make([]*WarrantData, len(pb))
	for i, v := range pb {
		list[i] = warrantDataFromPB(v)
	}
	return list
}
