package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotstockfilter"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotStockFilter = 3215 //Qot_StockFilter	获取条件选股

func init() {
	workers[ProtoIDQotStockFilter] = protocol.NewGetter()
}

// 条件选股
func (api *FutuAPI) GetStockFilter(ctx context.Context, market qotcommon.QotMarket, begin int32, num int32,
	filter *StockFilter) (*qotstockfilter.S2C, error) {

	if market == qotcommon.QotMarket_QotMarket_Unknown {
		return nil, ErrParameters
	}
	req := &qotstockfilter.Request{
		C2S: &qotstockfilter.C2S{
			Begin:  proto.Int32(begin),
			Num:    proto.Int32(num),
			Market: proto.Int32(int32(market)),
		},
	}
	filter.Filter(req.C2S)

	ch := make(chan *qotstockfilter.Response)
	if err := api.proto.RegisterGet(ProtoIDQotStockFilter, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C(), protocol.Error(resp)
	}
}

type StockFilter struct {
	Plate                     *qotcommon.Security                     // 板块
	BaseFilterList            []*qotstockfilter.BaseFilter            // 简单指标过滤器
	AccumulateFilterList      []*qotstockfilter.AccumulateFilter      // 累积指标过滤器
	FinancialFilterList       []*qotstockfilter.FinancialFilter       // 财务指标过滤器
	PatternFilterList         []*qotstockfilter.PatternFilter         // 形态技术指标过滤器
	CustomIndicatorFilterList []*qotstockfilter.CustomIndicatorFilter // 自定义技术指标过滤器
}

func (filter *StockFilter) Filter(req *qotstockfilter.C2S) {
	if filter != nil {
		req.Plate = filter.Plate
		if len(filter.BaseFilterList) != 0 {
			req.BaseFilterList = filter.BaseFilterList
		}
		if len(filter.AccumulateFilterList) != 0 {
			req.AccumulateFilterList = filter.AccumulateFilterList
		}
		if len(filter.FinancialFilterList) != 0 {
			req.FinancialFilterList = filter.FinancialFilterList
		}
		if len(filter.PatternFilterList) != 0 {
			req.PatternFilterList = filter.PatternFilterList
		}
		if len(filter.CustomIndicatorFilterList) != 0 {
			req.CustomIndicatorFilterList = filter.CustomIndicatorFilterList
		}
	}
}
