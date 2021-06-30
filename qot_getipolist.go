package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetipolist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetIpoList = 3217 //Qot_GetIpoList	获取新股
)

// 获取 IPO 信息
func (api *FutuAPI) GetIPOList(ctx context.Context, market qotcommon.QotMarket) ([]*IPOData, error) {
	req := qotgetipolist.Request{
		C2S: &qotgetipolist.C2S{
			Market: (*int32)(&market),
		},
	}
	ch := make(qotgetipolist.ResponseChan)
	if err := api.get(ProtoIDQotGetIpoList, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return ipoDataListFromPB(resp.GetS2C().GetIpoList()), protocol.Error(resp)
	}
}

// 新股 IPO 数据
type IPOData struct {
	Basic    *BasicIPOData // IPO 基本数据
	CNExData *CNIPOExData  // A 股 IPO 额外数据
	HKExData *HKIPOExData  // 港股 IPO 额外数据
	USExData *USIPOExData  // 美股 IPO 额外数据
}

func ipoDataFromPB(pb *qotgetipolist.IpoData) *IPOData {
	if pb == nil {
		return nil
	}
	return &IPOData{
		Basic:    basicIPODataFromPB(pb.GetBasic()),
		CNExData: cnIPOExDataFromPB(pb.GetCnExData()),
		HKExData: hkIPOExDataFromPB(pb.GetHkExData()),
		USExData: usIPOExDataFromPB(pb.GetUsExData()),
	}
}

func ipoDataListFromPB(pb []*qotgetipolist.IpoData) []*IPOData {
	if pb == nil {
		return nil
	}
	list := make([]*IPOData, len(pb))
	for i, v := range pb {
		list[i] = ipoDataFromPB(v)
	}
	return list
}

// IPO 基本数据
type BasicIPOData struct {
	Security      *Security // Qot_Common::QotMarket 股票市场，支持沪股和深股，且沪股和深股不做区分都代表 A 股市场。
	Name          string    // 股票名称
	ListTime      string    // 上市日期字符串
	ListTimestamp float64   // 上市日期时间戳
}

func basicIPODataFromPB(pb *qotgetipolist.BasicIpoData) *BasicIPOData {
	if pb == nil {
		return nil
	}
	return &BasicIPOData{
		Security:      securityFromPB(pb.GetSecurity()),
		Name:          pb.GetName(),
		ListTime:      pb.GetListTime(),
		ListTimestamp: pb.GetListTimestamp(),
	}
}

type CNIPOExData struct {
	ApplyCode              string            // 申购代码
	IssueSize              int64             // 发行总数
	OnlineIssueSize        int64             // 网上发行量
	ApplyUpperLimit        int64             // 申购上限
	ApplyLimitMarketValue  int64             // 顶格申购需配市值
	IsEstimateIPOPrice     bool              // 是否预估发行价
	IPOPrice               float64           // 发行价 预估值会因为募集资金、发行数量、发行费用等数据变动而变动，仅供参考。实际数据公布后会第一时间更新。
	IndustryPERate         float64           // 行业市盈率
	IsEstimateWinningRatio bool              // 是否预估中签率
	WinningRatio           float64           // 中签率 该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。预估值会因为募集资金、发行数量、发行费用等数据变动而变动，仅供参考。实际数据公布后会第一时间更新。
	IssuePERate            float64           // 发行市盈率
	ApplyTime              string            // 申购日期字符串
	ApplyTimestamp         float64           // 申购日期时间戳
	WinningTime            string            // 公布中签日期字符串
	WinningTimestamp       float64           // 公布中签日期时间戳
	IsHasWon               bool              // 是否已经公布中签号
	WinningNumData         []*WinningNumData // Qot_GetIpoList::WinningNumData 中签号数据，对应 PC 中"公布中签日期的已公布"
}

func cnIPOExDataFromPB(pb *qotgetipolist.CNIpoExData) *CNIPOExData {
	if pb == nil {
		return nil
	}
	return &CNIPOExData{
		ApplyCode:              pb.GetApplyCode(),
		IssueSize:              pb.GetIssueSize(),
		OnlineIssueSize:        pb.GetOnlineIssueSize(),
		ApplyUpperLimit:        pb.GetApplyUpperLimit(),
		ApplyLimitMarketValue:  pb.GetApplyLimitMarketValue(),
		IsEstimateIPOPrice:     pb.GetIsEstimateIpoPrice(),
		IPOPrice:               pb.GetIpoPrice(),
		IndustryPERate:         pb.GetIndustryPeRate(),
		IsEstimateWinningRatio: pb.GetIsEstimateWinningRatio(),
		WinningRatio:           pb.GetWinningRatio(),
		IssuePERate:            pb.GetIssuePeRate(),
		ApplyTime:              pb.GetApplyTime(),
		ApplyTimestamp:         pb.GetApplyTimestamp(),
		WinningTime:            pb.GetWinningTime(),
		WinningTimestamp:       pb.GetWinningTimestamp(),
		IsHasWon:               pb.GetIsHasWon(),
		WinningNumData:         winningNumDataListFromPB(pb.GetWinningNumData()),
	}
}

type WinningNumData struct {
	WinningName string // 分组名
	WinningInfo string // 中签号信息
}

func winningNumDataFromPB(pb *qotgetipolist.WinningNumData) *WinningNumData {
	if pb == nil {
		return nil
	}
	return &WinningNumData{
		WinningName: pb.GetWinningName(),
		WinningInfo: pb.GetWinningInfo(),
	}
}

func winningNumDataListFromPB(pb []*qotgetipolist.WinningNumData) []*WinningNumData {
	if pb == nil {
		return nil
	}
	list := make([]*WinningNumData, len(pb))
	for i, v := range pb {
		list[i] = winningNumDataFromPB(v)
	}
	return list
}

type HKIPOExData struct {
	IPOPriceMin       float64 // 最低发售价
	IPOPriceMax       float64 // 最高发售价
	ListPrice         float64 // 上市价
	LotSize           int32   // 每手股数
	EntrancePrice     float64 // 入场费
	IsSubscribeStatus bool    // 是否为认购状态，True-认购中，False-待上市
	ApplyEndTime      string  // 截止认购日期字符串
	ApplyEndTimestamp float64 // 截止认购日期时间戳 因需处理认购手续，富途认购截止时间会早于交易所公布的日期。
}

func hkIPOExDataFromPB(pb *qotgetipolist.HKIpoExData) *HKIPOExData {
	if pb == nil {
		return nil
	}
	return &HKIPOExData{
		IPOPriceMin:       pb.GetIpoPriceMin(),
		IPOPriceMax:       pb.GetIpoPriceMax(),
		ListPrice:         pb.GetListPrice(),
		LotSize:           pb.GetLotSize(),
		EntrancePrice:     pb.GetEntrancePrice(),
		IsSubscribeStatus: pb.GetIsSubscribeStatus(),
		ApplyEndTime:      pb.GetApplyEndTime(),
		ApplyEndTimestamp: pb.GetApplyEndTimestamp(),
	}
}

type USIPOExData struct {
	IPOPriceMin float64 // 最低发行价
	IPOPriceMax float64 // 最高发行价
	IssueSize   int64   // 发行量
}

func usIPOExDataFromPB(pb *qotgetipolist.USIpoExData) *USIPOExData {
	if pb == nil {
		return nil
	}
	return &USIPOExData{
		IPOPriceMin: pb.GetIpoPriceMin(),
		IPOPriceMax: pb.GetIpoPriceMax(),
		IssueSize:   pb.GetIssueSize(),
	}
}
