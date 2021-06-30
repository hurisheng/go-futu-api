package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetsecuritysnapshot"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetSecuritySnapshot = 3203 //Qot_GetSecuritySnapshot	获取股票快照
)

// 获取快照
func (api *FutuAPI) GetMarketSnapshot(ctx context.Context, securities []*Security) ([]*Snapshot, error) {
	// 请求参数
	req := qotgetsecuritysnapshot.Request{
		C2S: &qotgetsecuritysnapshot.C2S{
			SecurityList: securityList(securities).pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetsecuritysnapshot.ResponseChan)
	if err := api.get(ProtoIDQotGetSecuritySnapshot, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return snapshotListFromPB(resp.GetS2C().GetSnapshotList()), protocol.Error(resp)
	}
}

type Snapshot struct {
	Basic         *SnapshotBasicData     //*快照基本数据
	EquityExData  *EquitySnapshotExData  //正股快照额外数据
	WarrantExData *WarrantSnapshotExData //窝轮快照额外数据
	OptionExData  *OptionSnapshotExData  //期权快照额外数据
	IndexExData   *IndexSnapshotExData   //指数快照额外数据
	PlateExData   *PlateSnapshotExData   //板块快照额外数据
	FutureExData  *FutureSnapshotExData  //期货类型额外数据
	TrustExData   *TrustSnapshotExData   //基金类型额外数据
}

func snapshotFromPB(pb *qotgetsecuritysnapshot.Snapshot) *Snapshot {
	if pb == nil {
		return nil
	}
	return &Snapshot{
		Basic:         snapshotBasicDataFromPB(pb.GetBasic()),
		EquityExData:  equetyExDataFromPB(pb.GetEquityExData()),
		WarrantExData: warrantExDataFromPB(pb.GetWarrantExData()),
		OptionExData:  optionExDataFromPB(pb.GetOptionExData()),
		IndexExData:   indexExDataFromPB(pb.GetIndexExData()),
		PlateExData:   plateExDataFromPB(pb.GetPlateExData()),
		FutureExData:  futureExDataFromPB(pb.GetFutureExData()),
		TrustExData:   trustExDataFromPB(pb.GetTrustExData()),
	}
}

func snapshotListFromPB(pb []*qotgetsecuritysnapshot.Snapshot) []*Snapshot {
	if pb == nil {
		return nil
	}
	s := make([]*Snapshot, len(pb))
	for i, v := range pb {
		s[i] = snapshotFromPB(v)
	}
	return s
}

type SnapshotBasicData struct {
	Security                *Security                //*证券
	Type                    qotcommon.SecurityType   //*Qot_Common.SecurityType，证券类型
	IsSuspend               bool                     //*是否停牌
	ListTime                string                   //*上市时间字符串
	LotSize                 int32                    //*每手数量
	PriceSpread             float64                  //*价差
	UpdateTime              string                   //*更新时间字符串
	HighPrice               float64                  //*最高价
	OpenPrice               float64                  //*开盘价
	LowPrice                float64                  //*最低价
	LastClosePrice          float64                  //*昨收价
	CurPrice                float64                  //*最新价
	Volume                  int64                    //*成交量
	Turnover                float64                  //*成交额
	TurnoverRate            float64                  //换手率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	ListTimestamp           float64                  //上市时间戳
	UpdateTimestamp         float64                  //更新时间戳
	AskPrice                float64                  //卖价
	BidPrice                float64                  //买价
	AskVol                  int64                    //卖量
	BidVol                  int64                    //买量
	EnableMargin            bool                     // 是否可融资，如果为 true，后两个字段才有意义
	MortgageRatio           float64                  // 股票抵押率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	LongMarginInitialRatio  float64                  // 融资初始保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	EnableShortSell         bool                     // 是否可卖空，如果为 true，后三个字段才有意义
	ShortSellRate           float64                  // 卖空参考利率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	ShortAvailableVolume    int64                    // 剩余可卖空数量（股）
	ShortMarginInitialRatio float64                  // 卖空（融券）初始保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	Amplitude               float64                  // 振幅（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	AvgPrice                float64                  // 平均价
	BidAskRatio             float64                  // 委比（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	VolumeRatio             float64                  // 量比
	Highest52WeeksPrice     float64                  // 52周最高价
	Lowest52WeeksPrice      float64                  // 52周最低价
	HighestHistoryPrice     float64                  // 历史最高价
	LowestHistoryPrice      float64                  // 历史最低价
	PreMarket               *PreAfterMarketData      //Qot_Common::PreAfterMarketData 盘前数据
	AfterMarket             *PreAfterMarketData      //Qot_Common::PreAfterMarketData 盘后数据
	SecStatus               qotcommon.SecurityStatus //Qot_Common::SecurityStatus 股票状态
	ClosePrice5Minute       float64                  //5分钟收盘价
}

func snapshotBasicDataFromPB(pb *qotgetsecuritysnapshot.SnapshotBasicData) *SnapshotBasicData {
	if pb == nil {
		return nil
	}
	return &SnapshotBasicData{
		Security:                securityFromPB(pb.GetSecurity()),
		Type:                    qotcommon.SecurityType(pb.GetType()),
		IsSuspend:               pb.GetIsSuspend(),
		ListTime:                pb.GetListTime(),
		LotSize:                 pb.GetLotSize(),
		PriceSpread:             pb.GetPriceSpread(),
		UpdateTime:              pb.GetUpdateTime(),
		HighPrice:               pb.GetHighPrice(),
		OpenPrice:               pb.GetOpenPrice(),
		LowPrice:                pb.GetLowPrice(),
		LastClosePrice:          pb.GetLastClosePrice(),
		CurPrice:                pb.GetCurPrice(),
		Volume:                  pb.GetVolume(),
		Turnover:                pb.GetTurnover(),
		TurnoverRate:            pb.GetTurnoverRate(),
		ListTimestamp:           pb.GetListTimestamp(),
		UpdateTimestamp:         pb.GetUpdateTimestamp(),
		AskPrice:                pb.GetAskPrice(),
		BidPrice:                pb.GetBidPrice(),
		AskVol:                  pb.GetAskVol(),
		BidVol:                  pb.GetBidVol(),
		EnableMargin:            pb.GetEnableMargin(),
		MortgageRatio:           pb.GetMortgageRatio(),
		LongMarginInitialRatio:  pb.GetLongMarginInitialRatio(),
		EnableShortSell:         pb.GetEnableShortSell(),
		ShortSellRate:           pb.GetShortSellRate(),
		ShortAvailableVolume:    pb.GetShortAvailableVolume(),
		ShortMarginInitialRatio: pb.GetShortMarginInitialRatio(),
		Amplitude:               pb.GetAmplitude(),
		AvgPrice:                pb.GetAvgPrice(),
		BidAskRatio:             pb.GetBidAskRatio(),
		VolumeRatio:             pb.GetVolumeRatio(),
		Highest52WeeksPrice:     pb.GetHighest52WeeksPrice(),
		Lowest52WeeksPrice:      pb.GetLowest52WeeksPrice(),
		PreMarket:               preAfterMarketDataFromPB(pb.GetPreMarket()),
		AfterMarket:             preAfterMarketDataFromPB(pb.GetAfterMarket()),
		SecStatus:               qotcommon.SecurityStatus(pb.GetSecStatus()),
		ClosePrice5Minute:       pb.GetClosePrice5Minute(),
	}
}

type EquitySnapshotExData struct {
	IssuedShares         int64   // *发行股本，即总股本
	IssueMarketVal       float64 // *总市值 =总股本*当前价格
	NetAsset             float64 // *资产净值
	NetProfit            float64 // *盈利（亏损）
	EarningPershare      float64 // *每股盈利
	OutstandingShares    int64   // *流通股本
	OutstandingMarketVal float64 // *流通市值 =流通股本*当前价格
	NetAssetPershare     float64 // *每股净资产
	EYRate               float64 // *收益率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	PERate               float64 // *市盈率
	PBRate               float64 // *市净率
	PETTMRate            float64 // *市盈率 TTM
	DividendTTM          float64 // 股息 TTM，派息
	DividendRatioTTM     float64 // 股息率 TTM（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	DividendLFY          float64 // 股息 LFY，上一年度派息
	DividendLFYRatio     float64 // 股息率 LFY（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
}

func equetyExDataFromPB(pb *qotgetsecuritysnapshot.EquitySnapshotExData) *EquitySnapshotExData {
	if pb == nil {
		return nil
	}
	return &EquitySnapshotExData{
		IssuedShares:         pb.GetIssuedShares(),
		IssueMarketVal:       pb.GetIssuedMarketVal(),
		NetAsset:             pb.GetNetAsset(),
		NetProfit:            pb.GetNetProfit(),
		EarningPershare:      pb.GetEarningsPershare(),
		OutstandingShares:    pb.GetOutstandingShares(),
		OutstandingMarketVal: pb.GetOutstandingMarketVal(),
		NetAssetPershare:     pb.GetNetAssetPershare(),
		EYRate:               pb.GetEyRate(),
		PERate:               pb.GetPeRate(),
		PBRate:               pb.GetPbRate(),
		PETTMRate:            pb.GetPeTTMRate(),
		DividendTTM:          pb.GetDividendTTM(),
		DividendRatioTTM:     pb.GetDividendRatioTTM(),
		DividendLFY:          pb.GetDividendLFY(),
		DividendLFYRatio:     pb.GetDividendLFYRatio(),
	}
}

type WarrantSnapshotExData struct {
	ConversionRate     float64               //*换股比率
	WarrantType        qotcommon.WarrantType //*Qot_Common.WarrantType，窝轮类型
	StrikePrice        float64               //*行使价
	MaturityTime       string                //*到期日时间字符串
	EndTradeTime       string                //*最后交易日时间字符串
	Owner              *Security             //*所属正股
	RecoveryPrice      float64               //*收回价，仅牛熊证支持该字段
	StreetVolumn       int64                 //*街货量
	IssueVolumn        int64                 //*发行量
	StreetRate         float64               //*街货占比（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	Delta              float64               //*对冲值，仅认购认沽支持该字段
	ImpliedVolatility  float64               //*引申波幅，仅认购认沽支持该字段
	Premium            float64               //*溢价（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	MaturityTimestamp  float64               //到期日时间戳
	EndTradeTimestamp  float64               //最后交易日时间戳
	Leverage           float64               // 杠杆比率（倍）
	IPOP               float64               // 价内/价外（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	BreakEvenPoint     float64               // 打和点
	ConversionPrice    float64               // 换股价
	PriceRecoveryRatio float64               // 正股距收回价（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	Score              float64               // 综合评分
	UpperStrikePrice   float64               //上限价，仅界内证支持该字段
	LowerStrikePrice   float64               //下限价，仅界内证支持该字段
	InLinePriceStatus  qotcommon.PriceType   //Qot_Common.PriceType，界内界外，仅界内证支持该字段
	IssuerCode         string                //发行人代码
}

func warrantExDataFromPB(pb *qotgetsecuritysnapshot.WarrantSnapshotExData) *WarrantSnapshotExData {
	if pb == nil {
		return nil
	}
	return &WarrantSnapshotExData{
		ConversionRate:     pb.GetConversionRate(),
		WarrantType:        qotcommon.WarrantType(pb.GetWarrantType()),
		StrikePrice:        pb.GetStrikePrice(),
		MaturityTime:       pb.GetMaturityTime(),
		EndTradeTime:       pb.GetEndTradeTime(),
		Owner:              securityFromPB(pb.GetOwner()),
		RecoveryPrice:      pb.GetRecoveryPrice(),
		StreetVolumn:       pb.GetStreetVolumn(),
		IssueVolumn:        pb.GetIssueVolumn(),
		StreetRate:         pb.GetStreetRate(),
		Delta:              pb.GetDelta(),
		ImpliedVolatility:  pb.GetImpliedVolatility(),
		Premium:            pb.GetPremium(),
		MaturityTimestamp:  pb.GetMaturityTimestamp(),
		EndTradeTimestamp:  pb.GetEndTradeTimestamp(),
		Leverage:           pb.GetLeverage(),
		IPOP:               pb.GetIpop(),
		BreakEvenPoint:     pb.GetBreakEvenPoint(),
		ConversionPrice:    pb.GetConversionPrice(),
		PriceRecoveryRatio: pb.GetPriceRecoveryRatio(),
		Score:              pb.GetScore(),
		UpperStrikePrice:   pb.GetUpperStrikePrice(),
		LowerStrikePrice:   pb.GetLowerStrikePrice(),
		InLinePriceStatus:  qotcommon.PriceType(pb.GetInLinePriceStatus()),
		IssuerCode:         pb.GetIssuerCode(),
	}
}

type OptionSnapshotExData struct {
	Type                 qotcommon.OptionType      //*Qot_Common.OptionType，期权类型（按方向）
	Owner                *Security                 //*标的股
	StrikeTime           string                    //*行权日时间字符串
	StrikePrice          float64                   //*行权价
	ContractSize         int32                     //*每份合约数(整型数据)
	ContractSizeFloat    float64                   //*每份合约数（浮点型数据）
	OpenInterest         int32                     //*未平仓合约数
	ImpliedVolatility    float64                   //*隐含波动率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	Premium              float64                   //*溢价（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	Delta                float64                   //*希腊值 Delta
	Gamma                float64                   //*希腊值 Gamma
	Vega                 float64                   //*希腊值 Vega
	Theta                float64                   //*希腊值 Theta
	Rho                  float64                   //*希腊值 Rho
	StrikeTimestamp      float64                   //行权日时间戳
	IndexOptionType      qotcommon.IndexOptionType //Qot_Common.IndexOptionType，指数期权类型
	NetOpenInterest      int32                     //净未平仓合约数，仅港股期权适用
	ExpiryDateDistance   int32                     //距离到期日天数，负数表示已过期
	ContractNominalValue float64                   //合约名义金额，仅港股期权适用
	OwnerLotMultiplier   float64                   //相等正股手数，指数期权无该字段，仅港股期权适用
	OptionAreaType       qotcommon.OptionAreaType  //Qot_Common.OptionAreaType，期权类型（按行权时间）
	ContractMultiplier   float64                   //合约乘数
}

func optionExDataFromPB(pb *qotgetsecuritysnapshot.OptionSnapshotExData) *OptionSnapshotExData {
	if pb == nil {
		return nil
	}
	return &OptionSnapshotExData{
		Type:                 qotcommon.OptionType(pb.GetType()),
		Owner:                securityFromPB(pb.GetOwner()),
		StrikeTime:           pb.GetStrikeTime(),
		StrikePrice:          pb.GetStrikePrice(),
		ContractSize:         pb.GetContractSize(),
		ContractSizeFloat:    pb.GetContractSizeFloat(),
		OpenInterest:         pb.GetOpenInterest(),
		ImpliedVolatility:    pb.GetImpliedVolatility(),
		Premium:              pb.GetPremium(),
		Delta:                pb.GetDelta(),
		Gamma:                pb.GetGamma(),
		Vega:                 pb.GetVega(),
		Theta:                pb.GetTheta(),
		Rho:                  pb.GetRho(),
		StrikeTimestamp:      pb.GetStrikeTimestamp(),
		IndexOptionType:      qotcommon.IndexOptionType(pb.GetIndexOptionType()),
		NetOpenInterest:      pb.GetNetOpenInterest(),
		ExpiryDateDistance:   pb.GetExpiryDateDistance(),
		ContractNominalValue: pb.GetContractNominalValue(),
		OwnerLotMultiplier:   pb.GetOwnerLotMultiplier(),
		OptionAreaType:       qotcommon.OptionAreaType(pb.GetOptionAreaType()),
		ContractMultiplier:   pb.GetContractMultiplier(),
	}
}

type IndexSnapshotExData struct {
	RaiseCount int32 // 上涨支数
	FallCount  int32 // 下跌支数
	EqualCount int32 // 平盘支数
}

func indexExDataFromPB(pb *qotgetsecuritysnapshot.IndexSnapshotExData) *IndexSnapshotExData {
	if pb == nil {
		return nil
	}
	return &IndexSnapshotExData{
		RaiseCount: pb.GetRaiseCount(),
		FallCount:  pb.GetFallCount(),
		EqualCount: pb.GetEqualCount(),
	}
}

type PlateSnapshotExData struct {
	RaiseCount int32 // 上涨支数
	FallCount  int32 // 下跌支数
	EqualCount int32 // 平盘支数
}

func plateExDataFromPB(pb *qotgetsecuritysnapshot.PlateSnapshotExData) *PlateSnapshotExData {
	if pb == nil {
		return nil
	}
	return &PlateSnapshotExData{
		RaiseCount: pb.GetRaiseCount(),
		FallCount:  pb.GetFallCount(),
		EqualCount: pb.GetEqualCount(),
	}
}

type FutureSnapshotExData struct {
	LastSettlePrice    float64 //昨结
	Position           int32   //持仓量
	PositionChange     int32   //日增仓
	LastTradeTime      string  //最后交易日，只有非主连期货合约才有该字段
	LastTradeTimestamp float64 //最后交易日时间戳，只有非主连期货合约才有该字段
	IsMainContract     bool    //是否主连合约
}

func futureExDataFromPB(pb *qotgetsecuritysnapshot.FutureSnapshotExData) *FutureSnapshotExData {
	if pb == nil {
		return nil
	}
	return &FutureSnapshotExData{
		LastSettlePrice:    pb.GetLastSettlePrice(),
		Position:           pb.GetPosition(),
		PositionChange:     pb.GetPositionChange(),
		LastTradeTime:      pb.GetLastTradeTime(),
		LastTradeTimestamp: pb.GetLastTradeTimestamp(),
		IsMainContract:     pb.GetIsMainContract(),
	}
}

type TrustSnapshotExData struct {
	DividendYield    float64              //股息率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	Aum              float64              //资产规模
	OutstandingUnits int64                //总发行量
	NetAssetValue    float64              //单位净值
	Premium          float64              //溢价（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	AssetClass       qotcommon.AssetClass //Qot_Common.AssetClass，资产类别
}

func trustExDataFromPB(pb *qotgetsecuritysnapshot.TrustSnapshotExData) *TrustSnapshotExData {
	if pb == nil {
		return nil
	}
	return &TrustSnapshotExData{
		DividendYield:    pb.GetDividendYield(),
		Aum:              pb.GetAum(),
		OutstandingUnits: pb.GetOutstandingUnits(),
		NetAssetValue:    pb.GetNetAssetValue(),
		Premium:          pb.GetPremium(),
		AssetClass:       qotcommon.AssetClass(pb.GetAssetClass()),
	}
}
