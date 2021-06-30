package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotstockfilter"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotStockFilter = 3215 //Qot_StockFilter	获取条件选股
)

// 条件选股
func (api *FutuAPI) GetStockFilter(ctx context.Context, market qotcommon.QotMarket, begin int32, num int32, filter *StockFilter) (*StockFilterResult, error) {
	req := qotstockfilter.Request{
		C2S: &qotstockfilter.C2S{
			Begin:  &begin,
			Num:    &num,
			Market: (*int32)(&market),
		},
	}
	if filter != nil {
		req.C2S.Plate = filter.Plate.pb()
		req.C2S.BaseFilterList = baseFilterList(filter.BaseFilterList).pb()
		req.C2S.AccumulateFilterList = accumulateFilterList(filter.AccumulateFilterList).pb()
		req.C2S.FinancialFilterList = financialFilterList(filter.FinancialFilterList).pb()
	}
	ch := make(qotstockfilter.ResponseChan)
	if err := api.get(ProtoIDQotStockFilter, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return stockFilterResultFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type StockFilter struct {
	Plate                *Security           // 板块
	BaseFilterList       []*BaseFilter       // 简单指标过滤器
	AccumulateFilterList []*AccumulateFilter // 累积指标过滤器
	FinancialFilterList  []*FinancialFilter  // 财务指标过滤器
}

// 简单属性筛选
type BaseFilter struct {
	FieldName  qotstockfilter.StockField // StockField 简单属性
	FilterMin  *FilterDouble             // 区间下限（闭区间），不传代表下限为 -∞
	FilterMax  *FilterDouble             // 区间上限（闭区间），不传代表上限为 +∞
	IsNoFilter bool                      // 该字段是否不需要筛选，True：不筛选，False：筛选。不传默认不筛选
	SortDir    qotstockfilter.SortDir    // SortDir 排序方向，默认不排序。
}

func (f *BaseFilter) pb() *qotstockfilter.BaseFilter {
	if f == nil {
		return nil
	}
	return &qotstockfilter.BaseFilter{
		FieldName:  (*int32)(&f.FieldName),
		FilterMin:  f.FilterMin.pb(),
		FilterMax:  f.FilterMax.pb(),
		IsNoFilter: &f.IsNoFilter,
		SortDir:    (*int32)(&f.SortDir),
	}
}

type baseFilterList []*BaseFilter

func (f baseFilterList) pb() []*qotstockfilter.BaseFilter {
	if f == nil {
		return nil
	}
	list := make([]*qotstockfilter.BaseFilter, len(f))
	for i, v := range f {
		list[i] = v.pb()
	}
	return list
}

// 累积属性筛选
type AccumulateFilter struct {
	FieldName  qotstockfilter.AccumulateField // AccumulateField 累积属性
	FilterMin  *FilterDouble                  // 区间下限（闭区间），不传代表下限为 -∞
	FilterMax  *FilterDouble                  // 区间上限（闭区间），不传代表上限为 +∞
	IsNoFilter bool                           // 该字段是否不需要筛选，True：不筛选，False：筛选。不传默认不筛选
	SortDir    qotstockfilter.SortDir         // SortDir 排序方向，默认不排序。
	Days       int32                          // 近几日，累积时间
}

func (f *AccumulateFilter) pb() *qotstockfilter.AccumulateFilter {
	if f == nil {
		return nil
	}
	return &qotstockfilter.AccumulateFilter{
		FieldName:  (*int32)(&f.FieldName),
		FilterMin:  f.FilterMin.pb(),
		FilterMax:  f.FilterMax.pb(),
		IsNoFilter: &f.IsNoFilter,
		SortDir:    (*int32)(&f.SortDir),
		Days:       &f.Days,
	}
}

type accumulateFilterList []*AccumulateFilter

func (f accumulateFilterList) pb() []*qotstockfilter.AccumulateFilter {
	if f == nil {
		return nil
	}
	list := make([]*qotstockfilter.AccumulateFilter, len(f))
	for i, v := range f {
		list[i] = v.pb()
	}
	return list
}

// 财务属性筛选
type FinancialFilter struct {
	FiledName  qotstockfilter.FinancialField   // FinancialField 财务属性
	FilterMin  *FilterDouble                   // 区间下限（闭区间），不传代表下限为 -∞
	FilterMax  *FilterDouble                   // 区间上限（闭区间），不传代表上限为 +∞
	IsNoFilter bool                            // 该字段是否不需要筛选，True：不筛选，False：筛选。不传默认不筛选
	SortDir    qotstockfilter.SortDir          // SortDir 排序方向，默认不排序。
	Quarter    qotstockfilter.FinancialQuarter // FinancialQuarter 财报累积时间
}

func (f *FinancialFilter) pb() *qotstockfilter.FinancialFilter {
	if f == nil {
		return nil
	}
	return &qotstockfilter.FinancialFilter{
		FieldName:  (*int32)(&f.FiledName),
		FilterMin:  f.FilterMin.pb(),
		FilterMax:  f.FilterMax.pb(),
		IsNoFilter: &f.IsNoFilter,
		SortDir:    (*int32)(&f.SortDir),
		Quarter:    (*int32)(&f.Quarter),
	}
}

type financialFilterList []*FinancialFilter

func (f financialFilterList) pb() []*qotstockfilter.FinancialFilter {
	if f == nil {
		return nil
	}
	list := make([]*qotstockfilter.FinancialFilter, len(f))
	for i, v := range f {
		list[i] = v.pb()
	}
	return list
}

type StockFilterResult struct {
	LastPage bool         // 是否最后一页了,false:非最后一页,还有窝轮记录未返回; true:已是最后一页
	AllCount int32        // 该条件请求所有数据的个数
	DataList []*StockData // 返回的股票数据列表
}

func stockFilterResultFromPB(pb *qotstockfilter.S2C) *StockFilterResult {
	if pb == nil {
		return nil
	}
	return &StockFilterResult{
		LastPage: pb.GetLastPage(),
		AllCount: pb.GetAllCount(),
		DataList: stockDataListFromPB(pb.GetDataList()),
	}
}

// 返回的股票数据
type StockData struct {
	Security           *Security         // 股票
	Name               string            // 股票名称
	BaseDataList       []*BaseData       // 筛选后的简单指标属性数据
	AccumulateDataList []*AccumulateData // 筛选后的累积指标属性数据
	FinancialDataList  []*FinancialData  // 筛选后的财务指标属性数据
}

func stockDataFromPB(pb *qotstockfilter.StockData) *StockData {
	if pb == nil {
		return nil
	}
	return &StockData{
		Security:           securityFromPB(pb.GetSecurity()),
		Name:               pb.GetName(),
		BaseDataList:       baseDataListFromPB(pb.GetBaseDataList()),
		AccumulateDataList: accumulateDataListFromPB(pb.GetAccumulateDataList()),
		FinancialDataList:  financialDataListFromPB(pb.GetFinancialDataList()),
	}
}

func stockDataListFromPB(pb []*qotstockfilter.StockData) []*StockData {
	if pb == nil {
		return nil
	}
	list := make([]*StockData, len(pb))
	for i, v := range pb {
		list[i] = stockDataFromPB(v)
	}
	return list
}

// 简单属性数据
type BaseData struct {
	FieldName qotstockfilter.StockField // StockField 简单属性
	Value     float64
}

func baseDataFromPB(pb *qotstockfilter.BaseData) *BaseData {
	if pb == nil {
		return nil
	}
	return &BaseData{
		FieldName: qotstockfilter.StockField(pb.GetFieldName()),
		Value:     pb.GetValue(),
	}
}

func baseDataListFromPB(pb []*qotstockfilter.BaseData) []*BaseData {
	if pb == nil {
		return nil
	}
	list := make([]*BaseData, len(pb))
	for i, v := range pb {
		list[i] = baseDataFromPB(v)
	}
	return list
}

// 累积属性数据
type AccumulateData struct {
	FieldName qotstockfilter.AccumulateField // AccumulateField 累积属性
	Value     float64
	Days      int32 // 近几日，累积时间
}

func accumulateDataFromPB(pb *qotstockfilter.AccumulateData) *AccumulateData {
	if pb == nil {
		return nil
	}
	return &AccumulateData{
		FieldName: qotstockfilter.AccumulateField(pb.GetFieldName()),
		Value:     pb.GetValue(),
		Days:      pb.GetDays(),
	}
}

func accumulateDataListFromPB(pb []*qotstockfilter.AccumulateData) []*AccumulateData {
	if pb == nil {
		return nil
	}
	list := make([]*AccumulateData, len(pb))
	for i, v := range pb {
		list[i] = accumulateDataFromPB(v)
	}
	return list
}

// 财务属性数据
type FinancialData struct {
	FieldName qotstockfilter.FinancialField // FinancialField 财务属性
	Value     float64
	Quarter   qotstockfilter.FinancialQuarter // FinancialQuarter 财报累积时间
}

func financialDataFromPB(pb *qotstockfilter.FinancialData) *FinancialData {
	if pb == nil {
		return nil
	}
	return &FinancialData{
		FieldName: qotstockfilter.FinancialField(pb.GetFieldName()),
		Value:     pb.GetValue(),
		Quarter:   qotstockfilter.FinancialQuarter(pb.GetQuarter()),
	}
}

func financialDataListFromPB(pb []*qotstockfilter.FinancialData) []*FinancialData {
	if pb == nil {
		return nil
	}
	list := make([]*FinancialData, len(pb))
	for i, v := range pb {
		list[i] = financialDataFromPB(v)
	}
	return list
}
