package futuapi

import "github.com/hurisheng/go-futu-api/pb/trdcommon"

type TrdAcc struct {
	TrdEnv            trdcommon.TrdEnv       //交易环境，参见 TrdEnv 的枚举定义
	AccID             uint64                 //业务账号
	TrdMarketAuthList []trdcommon.TrdMarket  //业务账户支持的交易市场权限，即此账户能交易那些市场, 可拥有多个交易市场权限，目前仅单个，取值参见 TrdMarket 的枚举定义
	AccType           trdcommon.TrdAccType   //账户类型，取值见 TrdAccType
	CardNum           string                 //卡号
	SecurityFirm      trdcommon.SecurityFirm //所属券商，取值见SecurityFirm
	SimAccType        trdcommon.SimAccType   //模拟交易账号类型，取值见SimAccType
}

func trdAccFromPB(pb *trdcommon.TrdAcc) *TrdAcc {
	if pb == nil {
		return nil
	}
	return &TrdAcc{
		TrdEnv:            trdcommon.TrdEnv(pb.GetTrdEnv()),
		AccID:             pb.GetAccID(),
		TrdMarketAuthList: trdMarketListFromPB(pb.GetTrdMarketAuthList()),
		AccType:           trdcommon.TrdAccType(pb.GetAccType()),
		CardNum:           pb.GetCardNum(),
		SecurityFirm:      trdcommon.SecurityFirm(pb.GetSecurityFirm()),
		SimAccType:        trdcommon.SimAccType(pb.GetSimAccType()),
	}
}

func trdAccListFromPB(pb []*trdcommon.TrdAcc) []*TrdAcc {
	if pb == nil {
		return nil
	}
	list := make([]*TrdAcc, len(pb))
	for i, v := range pb {
		list[i] = trdAccFromPB(v)
	}
	return list
}

func trdMarketListFromPB(pb []int32) []trdcommon.TrdMarket {
	if pb == nil {
		return nil
	}
	list := make([]trdcommon.TrdMarket, len(pb))
	for i, v := range pb {
		list[i] = trdcommon.TrdMarket(v)
	}
	return list
}

// 交易协议公共参数头
type TrdHeader struct {
	TrdEnv    trdcommon.TrdEnv    //交易环境, 参见 TrdEnv 的枚举定义
	AccID     uint64              //业务账号, 业务账号与交易环境、市场权限需要匹配，否则会返回错误
	TrdMarket trdcommon.TrdMarket //交易市场, 参见 TrdMarket 的枚举定义
}

func (h *TrdHeader) pb() *trdcommon.TrdHeader {
	if h == nil {
		return nil
	}
	return &trdcommon.TrdHeader{
		TrdEnv:    (*int32)(&h.TrdEnv),
		AccID:     &h.AccID,
		TrdMarket: (*int32)(&h.TrdMarket),
	}
}

func trdHeaderFromPB(pb *trdcommon.TrdHeader) *TrdHeader {
	if pb == nil {
		return nil
	}
	return &TrdHeader{
		TrdEnv:    trdcommon.TrdEnv(pb.GetTrdEnv()),
		AccID:     pb.GetAccID(),
		TrdMarket: trdcommon.TrdMarket(pb.GetTrdMarket()),
	}
}

// 账户资金
type Funds struct {
	Power             float64 //最大购买力（做多）
	TotalAssets       float64 //资产净值
	Cash              float64 //现金
	MarketVal         float64 //证券市值, 仅证券账户适用
	FrozenCash        float64 //冻结资金
	DebtCash          float64 //计息金额
	AvlWithdrawalCash float64 //现金可提，仅证券账户适用

	Currency          trdcommon.Currency      //币种，本结构体资金相关的货币类型，取值参见 Currency，期货适用
	AvailableFunds    float64                 //可用资金，期货适用
	UnrealizedPL      float64                 //未实现盈亏，期货适用
	RealizedPL        float64                 //已实现盈亏，期货适用
	RiskLevel         trdcommon.CltRiskLevel  //风控状态，参见 CltRiskLevel, 期货适用
	InitialMargin     float64                 //初始保证金
	MaintenanceMargin float64                 //维持保证金
	CashInfoList      []*AccCashInfo          //分币种的现金信息，期货适用
	MaxPowerShort     float64                 //卖空购买力
	NetCashPower      float64                 //现金购买力
	LongMV            float64                 //多头市值
	ShortMV           float64                 //空头市值
	PendingAsset      float64                 //在途资产
	MaxWithdrawal     float64                 //融资可提，仅证券账户适用
	RiskStatus        trdcommon.CltRiskStatus //风险状态，参见 [CltRiskStatus]，证券账户适用，共分 9 个等级，LEVEL1是最安全，LEVEL9是最危险
	MarginCallMargin  float64                 //	Margin Call 保证金
}

func fundsFromPB(pb *trdcommon.Funds) *Funds {
	if pb == nil {
		return nil
	}
	return &Funds{
		Power:             pb.GetPower(),
		TotalAssets:       pb.GetTotalAssets(),
		Cash:              pb.GetCash(),
		MarketVal:         pb.GetMarketVal(),
		FrozenCash:        pb.GetFrozenCash(),
		DebtCash:          pb.GetDebtCash(),
		AvlWithdrawalCash: pb.GetAvlWithdrawalCash(),

		Currency:          trdcommon.Currency(pb.GetCurrency()),
		AvailableFunds:    pb.GetAvailableFunds(),
		UnrealizedPL:      pb.GetUnrealizedPL(),
		RealizedPL:        pb.GetRealizedPL(),
		RiskLevel:         trdcommon.CltRiskLevel(pb.GetRiskLevel()),
		InitialMargin:     pb.GetInitialMargin(),
		MaintenanceMargin: pb.GetMaintenanceMargin(),
		CashInfoList:      accCashInfoListFromPB(pb.GetCashInfoList()),
		MaxPowerShort:     pb.GetMaxPowerShort(),
		NetCashPower:      pb.GetNetCashPower(),
		LongMV:            pb.GetLongMv(),
		ShortMV:           pb.GetShortMv(),
		PendingAsset:      pb.GetPendingAsset(),
		MaxWithdrawal:     pb.GetMaxWithdrawal(),
		RiskStatus:        trdcommon.CltRiskStatus(pb.GetRiskStatus()),
		MarginCallMargin:  pb.GetMarginCallMargin(),
	}
}

type AccCashInfo struct {
	Currency         trdcommon.Currency // 货币类型，取值参考 Currency
	Cash             float64            // 现金结余
	AvailableBalance float64            // 现金可提金额
}

func accCashInfoFromPB(pb *trdcommon.AccCashInfo) *AccCashInfo {
	if pb == nil {
		return nil
	}
	return &AccCashInfo{
		Currency:         trdcommon.Currency(pb.GetCurrency()),
		Cash:             pb.GetCash(),
		AvailableBalance: pb.GetAvailableBalance(),
	}
}

func accCashInfoListFromPB(pb []*trdcommon.AccCashInfo) []*AccCashInfo {
	if pb == nil {
		return nil
	}
	list := make([]*AccCashInfo, len(pb))
	for i, v := range pb {
		list[i] = accCashInfoFromPB(v)
	}
	return list
}

// 最大可交易数量
type MaxTrdQtys struct {
	MaxCashBuy          float64 //不使用融资，仅自己的现金最大可买整手股数，期货此字段值为0
	MaxCashAndMarginBuy float64 //使用融资，自己的现金 + 融资资金总共的最大可买整手股数，期货不适用
	MaxPositionSell     float64 //不使用融券(卖空)，仅自己的持仓最大可卖整手股数
	MaxSellShort        float64 //使用融券(卖空)，最大可卖空整手股数，不包括多仓，期货不适用
	MaxBuyBack          float64 //卖空后，需要买回的最大整手股数。因为卖空后，必须先买回已卖空的股数，还掉股票，才能再继续买多。期货不适用
	LongRequiredIM      float64 //开多仓每张合约初始保证金。当前仅期货和期权适用（最低 FutuOpenD 版本要求：5.0.1310）
	ShortRequiredIM     float64 //开空仓每张合约初始保证金。当前仅期货和期权适用（最低 FutuOpenD 版本要求：5.0.1310）
}

func maxTrdQtysFromPB(pb *trdcommon.MaxTrdQtys) *MaxTrdQtys {
	if pb == nil {
		return nil
	}
	return &MaxTrdQtys{
		MaxCashBuy:          pb.GetMaxCashBuy(),
		MaxCashAndMarginBuy: pb.GetMaxCashAndMarginBuy(),
		MaxPositionSell:     pb.GetMaxPositionSell(),
		MaxSellShort:        pb.GetMaxSellShort(),
		MaxBuyBack:          pb.GetMaxBuyBack(),
		LongRequiredIM:      pb.GetLongRequiredIM(),
		ShortRequiredIM:     pb.GetShortRequiredIM(),
	}
}

type TrdFilterConditions struct {
	CodeList []string //代码过滤，只返回包含这些代码的数据，没传不过滤
	IDList   []uint64 //ID 主键过滤，只返回包含这些 ID 的数据，没传不过滤，订单是 orderID、成交是 fillID、持仓是 positionID
	Begin    string   //开始时间，严格按 YYYY-MM-DD HH:MM:SS 或 YYYY-MM-DD HH:MM:SS.MS 格式传，对持仓无效，拉历史数据必须填
	End      string   //结束时间，严格按 YYYY-MM-DD HH:MM:SS 或 YYYY-MM-DD HH:MM:SS.MS 格式传，对持仓无效，拉历史数据必须填
}

func (c *TrdFilterConditions) pb() *trdcommon.TrdFilterConditions {
	if c == nil {
		return nil
	}
	cond := trdcommon.TrdFilterConditions{
		CodeList: c.CodeList,
		IdList:   c.IDList,
	}
	if c.Begin != "" {
		cond.BeginTime = &c.Begin
	}
	if c.End != "" {
		cond.EndTime = &c.End
	}
	return &cond
}

type Position struct {
	PositionID   uint64                 //持仓 ID，一条持仓的唯一标识
	PositionSide trdcommon.PositionSide //持仓方向，参见 PositionSide 的枚举定
	Code         string                 //代码
	Name         string                 //名称
	Qty          float64                //持有数量，2位精度，期权单位是"张"，下同
	CanSellQty   float64                //可卖数量
	Price        float64                //市价，3位精度，期货为2位精度
	CostPrice    float64                //成本价，无精度限制，期货为2位精度，如果没传，代表此时此值无效,
	Val          float64                //市值，3位精度, 期货此字段值为0
	PLVal        float64                //盈亏金额，3位精度，期货为2位精度
	PLRatio      float64                //盈亏百分比(如 plRatio 等于8.8代表涨8.8%)，无精度限制，如果没传，代表此时此值无效
	SecMarket    trdcommon.TrdSecMarket //证券所属市场，参见 TrdSecMarket 的枚举定义
	//以下是此持仓今日统计
	TdPLVal      float64 //今日盈亏金额，3位精度，下同, 期货为2位精度
	TdTrdVal     float64 //今日交易额，期货不适用
	TdBuyVal     float64 //今日买入总额，期货不适用
	TdBuyQty     float64 //今日买入总量，期货不适用
	TdSellVal    float64 //今日卖出总额，期货不适用
	TdSellQty    float64 //今日卖出总量，期货不适用
	UnrealizedPL float64 //未实现盈亏（仅期货账户适用）
	RelizedPL    float64 //已实现盈亏（仅期货账户适用）
}

func positionFromPB(pb *trdcommon.Position) *Position {
	if pb == nil {
		return nil
	}
	return &Position{
		PositionID:   pb.GetPositionID(),
		PositionSide: trdcommon.PositionSide(pb.GetPositionSide()),
		Code:         pb.GetCode(),
		Name:         pb.GetName(),
		Qty:          pb.GetQty(),
		CanSellQty:   pb.GetCanSellQty(),
		Price:        pb.GetPrice(),
		CostPrice:    pb.GetCostPrice(),
		Val:          pb.GetVal(),
		PLVal:        pb.GetPlVal(),
		PLRatio:      pb.GetPlRatio(),
		SecMarket:    trdcommon.TrdSecMarket(pb.GetSecMarket()),
		TdPLVal:      pb.GetTdPlVal(),
		TdTrdVal:     pb.GetTdTrdVal(),
		TdBuyVal:     pb.GetTdBuyVal(),
		TdBuyQty:     pb.GetTdBuyQty(),
		TdSellVal:    pb.GetTdSellVal(),
		TdSellQty:    pb.GetTdSellQty(),
		UnrealizedPL: pb.GetUnrealizedPL(),
		RelizedPL:    pb.GetRealizedPL(),
	}
}

func positionListFromPB(pb []*trdcommon.Position) []*Position {
	if pb == nil {
		return nil
	}
	list := make([]*Position, len(pb))
	for i, v := range pb {
		list[i] = positionFromPB(v)
	}
	return list
}

type orderStatusList []trdcommon.OrderStatus

func (s orderStatusList) pb() []int32 {
	if s == nil {
		return nil
	}
	list := make([]int32, len(s))
	for i, v := range s {
		list[i] = int32(v)
	}
	return list
}

// 订单
type Order struct {
	TrdSide         trdcommon.TrdSide      //交易方向, 参见 TrdSide 的枚举定义
	OrderType       trdcommon.OrderType    //订单类型, 参见 OrderType 的枚举定义
	OrderStatus     trdcommon.OrderStatus  //订单状态, 参见 OrderStatus 的枚举定义
	OrderID         uint64                 //订单号
	OrderIDEx       string                 //扩展订单号(仅查问题时备用)
	Code            string                 //代码
	Name            string                 //名称
	Qty             float64                //订单数量，2位精度，期权单位是"张"
	Price           float64                //订单价格，3位精度
	CreateTime      string                 //创建时间，严格按 YYYY-MM-DD HH:MM:SS 或 YYYY-MM-DD HH:MM:SS.MS 格式传
	UpdateTime      string                 //最后更新时间，严格按 YYYY-MM-DD HH:MM:SS 或 YYYY-MM-DD HH:MM:SS.MS 格式传
	FillQty         float64                //成交数量，2位精度，期权单位是"张"
	FillAvgPrice    float64                //成交均价，无精度限制
	LastErrMsg      string                 //最后的错误描述，如果有错误，会有此描述最后一次错误的原因，无错误为空
	SecMarket       trdcommon.TrdSecMarket //证券所属市场，参见 TrdSecMarket 的枚举定义
	CreateTimestamp float64                //创建时间戳
	UpdateTimestamp float64                //最后更新时间戳
	Remark          string                 //用户备注字符串，最大长度64字节
}

func orderFromPB(pb *trdcommon.Order) *Order {
	if pb == nil {
		return nil
	}
	return &Order{
		TrdSide:         trdcommon.TrdSide(pb.GetTrdSide()),
		OrderType:       trdcommon.OrderType(pb.GetOrderType()),
		OrderStatus:     trdcommon.OrderStatus(pb.GetOrderStatus()),
		OrderID:         pb.GetOrderID(),
		OrderIDEx:       pb.GetOrderIDEx(),
		Code:            pb.GetCode(),
		Name:            pb.GetName(),
		Qty:             pb.GetQty(),
		Price:           pb.GetPrice(),
		CreateTime:      pb.GetCreateTime(),
		UpdateTime:      pb.GetUpdateTime(),
		FillQty:         pb.GetFillQty(),
		FillAvgPrice:    pb.GetFillAvgPrice(),
		LastErrMsg:      pb.GetLastErrMsg(),
		SecMarket:       trdcommon.TrdSecMarket(pb.GetSecMarket()),
		CreateTimestamp: pb.GetCreateTimestamp(),
		UpdateTimestamp: pb.GetUpdateTimestamp(),
		Remark:          pb.GetRemark(),
	}
}

func orderListFromPB(pb []*trdcommon.Order) []*Order {
	if pb == nil {
		return nil
	}
	list := make([]*Order, len(pb))
	for i, v := range pb {
		list[i] = orderFromPB(v)
	}
	return list
}

type OrderFill struct {
	TrdSide           trdcommon.TrdSide         //交易方向, 参见 TrdSide 的枚举定义
	FillID            uint64                    //成交号
	FillIDEx          string                    //扩展成交号(仅查问题时备用)
	OrderID           uint64                    //订单号
	OrderIDEx         string                    //扩展订单号(仅查问题时备用)
	Code              string                    //代码
	Name              string                    //名称
	Qty               float64                   //成交数量，2位精度，期权单位是"张"
	Price             float64                   //成交价格，3位精度
	CreateTime        string                    //创建时间（成交时间），严格按 YYYY-MM-DD HH:MM:SS 或 YYYY-MM-DD HH:MM:SS.MS 格式传
	CounterBrokerID   int32                     //对手经纪号，港股有效
	CounterBrokerName string                    //对手经纪名称，港股有效
	SecMarket         trdcommon.TrdSecMarket    //证券所属市场，参见 TrdSecMarket 的枚举定义
	CreateTimestamp   float64                   //创建时间戳
	UpdateTimestamp   float64                   //最后更新时间戳
	Status            trdcommon.OrderFillStatus //成交状态, 参见 OrderFillStatus 的枚举定义
}

func orderFillFromPB(pb *trdcommon.OrderFill) *OrderFill {
	if pb == nil {
		return nil
	}
	return &OrderFill{
		TrdSide:           trdcommon.TrdSide(pb.GetTrdSide()),
		FillID:            pb.GetFillID(),
		FillIDEx:          pb.GetFillIDEx(),
		OrderID:           pb.GetOrderID(),
		OrderIDEx:         pb.GetOrderIDEx(),
		Code:              pb.GetCode(),
		Name:              pb.GetName(),
		Qty:               pb.GetQty(),
		Price:             pb.GetPrice(),
		CreateTime:        pb.GetCreateTime(),
		CounterBrokerID:   pb.GetCounterBrokerID(),
		CounterBrokerName: pb.GetCounterBrokerName(),
		SecMarket:         trdcommon.TrdSecMarket(pb.GetSecMarket()),
		CreateTimestamp:   pb.GetCreateTimestamp(),
		UpdateTimestamp:   pb.GetUpdateTimestamp(),
		Status:            trdcommon.OrderFillStatus(pb.GetStatus()),
	}
}

func orderFillListFromPB(pb []*trdcommon.OrderFill) []*OrderFill {
	if pb == nil {
		return nil
	}
	list := make([]*OrderFill, len(pb))
	for i, v := range pb {
		list[i] = orderFillFromPB(v)
	}
	return list
}
