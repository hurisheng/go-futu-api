package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetwarrant"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetWarrant = 3210 //Qot_GetWarrant	获取窝轮

func init() {
	workers[ProtoIDQotGetWarrant] = protocol.NewGetter()
}

// 筛选窝轮
func (api *FutuAPI) GetWarrant(ctx context.Context, begin int32, num int32, sortField qotcommon.SortField, ascend bool,
	filter *WarrantFilter) (*qotgetwarrant.S2C, error) {

	if sortField == qotcommon.SortField_SortField_Unknow {
		return nil, ErrParameters
	}
	req := &qotgetwarrant.Request{
		C2S: &qotgetwarrant.C2S{
			Begin:     proto.Int32(begin),
			Num:       proto.Int32(num),
			SortField: proto.Int32(int32(sortField)),
			Ascend:    proto.Bool(ascend),
		},
	}
	filter.Filter(req.C2S)

	ch := make(chan *qotgetwarrant.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetWarrant, req, protocol.NewProtobufChan(ch)); err != nil {
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

type WarrantFilter struct {
	Owner                 *qotcommon.Security     //所属正股
	TypeList              []qotcommon.WarrantType //Qot_Common.WarrantType，窝轮类型过滤列表
	IssuerList            []qotcommon.Issuer      //Qot_Common.Issuer，发行人过滤列表
	IpoPeriod             qotcommon.IpoPeriod     //Qot_Common.IpoPeriod，上市日
	PriceType             qotcommon.PriceType     //Qot_Common.PriceType，价内/价外（暂不支持界内证的界内外筛选）
	Status                qotcommon.WarrantStatus //Qot_Common.WarrantStatus，窝轮状态
	MaturityTimeMin       string                  //到期日，到期日范围的开始时间戳
	MaturityTimeMax       string                  //到期日范围的结束时间戳
	CurPriceMin           *OptionalDouble         //最新价的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	CurPriceMax           *OptionalDouble         //最新价的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	StrikePriceMin        *OptionalDouble         //行使价的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	StrikePriceMax        *OptionalDouble         //行使价的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	StreetMin             *OptionalDouble         //街货占比的过滤下限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	StreetMax             *OptionalDouble         //街货占比的过滤上限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	ConversionMin         *OptionalDouble         //换股比率的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	ConversionMax         *OptionalDouble         //换股比率的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	VolMin                *OptionalUInt64         //成交量的过滤下限（闭区间），不传代表下限为 -∞
	VolMax                *OptionalUInt64         //成交量的过滤上限（闭区间），不传代表上限为 +∞
	PremiumMin            *OptionalDouble         //溢价的过滤下限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	PremiumMax            *OptionalDouble         //溢价的过滤上限（闭区间），该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	LeverageRatioMin      *OptionalDouble         //杠杆比率的过滤下限（闭区间），不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	LeverageRatioMax      *OptionalDouble         //杠杆比率的过滤上限（闭区间），不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	DeltaMin              *OptionalDouble         //对冲值的过滤下限（闭区间），仅认购认沽支持此字段过滤，不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	DeltaMax              *OptionalDouble         //对冲值的过滤上限（闭区间），仅认购认沽支持此字段过滤，不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	ImplieMin             *OptionalDouble         //引伸波幅的过滤下限（闭区间），仅认购认沽支持此字段过滤，不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	ImplieMax             *OptionalDouble         //引伸波幅的过滤上限（闭区间），仅认购认沽支持此字段过滤，不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	RecoveryPriceMin      *OptionalDouble         //收回价的过滤下限（闭区间），仅牛熊证支持此字段过滤，不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	RecoveryPriceMax      *OptionalDouble         //收回价的过滤上限（闭区间），仅牛熊证支持此字段过滤，不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
	PriceRecoveryRatioMin *OptionalDouble         //正股距收回价，的过滤下限（闭区间），仅牛熊证支持此字段过滤。该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表下限为 -∞（精确到小数点后 3 位，超出部分会被舍弃）
	PriceRecoveryRatioMax *OptionalDouble         //正股距收回价，的过滤上限（闭区间），仅牛熊证支持此字段过滤。该字段为百分比字段，默认不展示 %，如 20 实际对应 20%。不传代表上限为 +∞（精确到小数点后 3 位，超出部分会被舍弃）
}

func (filter *WarrantFilter) Filter(req *qotgetwarrant.C2S) {
	if filter != nil {
		if filter.Owner != nil {
			req.Owner = filter.Owner
		}
		if filter.TypeList != nil {
			req.TypeList = make([]int32, len(filter.TypeList))
			for i, v := range filter.TypeList {
				req.TypeList[i] = int32(v)
			}
		}
		if filter.IssuerList != nil {
			req.IssuerList = make([]int32, len(filter.IssuerList))
			for i, v := range filter.IssuerList {
				req.IssuerList[i] = int32(v)
			}
		}
		if filter.IpoPeriod != qotcommon.IpoPeriod_IpoPeriod_Unknow {
			req.IpoPeriod = proto.Int32(int32(filter.IpoPeriod))
		}
		if filter.PriceType != qotcommon.PriceType_PriceType_Unknow {
			req.PriceType = proto.Int32(int32(filter.PriceType))
		}
		if filter.Status != qotcommon.WarrantStatus_WarrantStatus_Unknow {
			req.Status = proto.Int32(int32(filter.Status))
		}
		if filter.MaturityTimeMin != "" {
			req.MaturityTimeMin = proto.String(filter.MaturityTimeMin)
		}
		if filter.MaturityTimeMax != "" {
			req.MaturityTimeMax = proto.String(filter.MaturityTimeMax)
		}
		if filter.CurPriceMin != nil {
			req.CurPriceMin = proto.Float64(filter.CurPriceMin.Value)
		}
		if filter.CurPriceMax != nil {
			req.CurPriceMax = proto.Float64(filter.CurPriceMax.Value)
		}
		if filter.StrikePriceMin != nil {
			req.StrikePriceMin = proto.Float64(filter.StrikePriceMin.Value)
		}
		if filter.StrikePriceMax != nil {
			req.StrikePriceMax = proto.Float64(filter.StrikePriceMax.Value)
		}
		if filter.StreetMin != nil {
			req.StreetMin = proto.Float64(filter.StreetMin.Value)
		}
		if filter.StreetMax != nil {
			req.StreetMax = proto.Float64(filter.StreetMax.Value)
		}
		if filter.ConversionMin != nil {
			req.ConversionMin = proto.Float64(filter.ConversionMin.Value)
		}
		if filter.ConversionMax != nil {
			req.ConversionMax = proto.Float64(filter.ConversionMax.Value)
		}
		if filter.VolMin != nil {
			req.VolMin = proto.Uint64(filter.VolMin.Value)
		}
		if filter.VolMax != nil {
			req.VolMax = proto.Uint64(filter.VolMax.Value)
		}
		if filter.VolMin != nil {
			req.VolMin = proto.Uint64(filter.VolMin.Value)
		}
		if filter.PremiumMin != nil {
			req.PremiumMin = proto.Float64(filter.PremiumMin.Value)
		}
		if filter.PremiumMax != nil {
			req.PremiumMax = proto.Float64(filter.PremiumMax.Value)
		}
		if filter.LeverageRatioMin != nil {
			req.LeverageRatioMin = proto.Float64(filter.LeverageRatioMin.Value)
		}
		if filter.LeverageRatioMax != nil {
			req.LeverageRatioMax = proto.Float64(filter.LeverageRatioMax.Value)
		}
		if filter.DeltaMin != nil {
			req.DeltaMin = proto.Float64(filter.DeltaMin.Value)
		}
		if filter.DeltaMax != nil {
			req.DeltaMax = proto.Float64(filter.DeltaMax.Value)
		}
		if filter.ImplieMin != nil {
			req.ImpliedMin = proto.Float64(filter.ImplieMin.Value)
		}
		if filter.ImplieMax != nil {
			req.ImpliedMax = proto.Float64(filter.ImplieMax.Value)
		}
		if filter.RecoveryPriceMin != nil {
			req.RecoveryPriceMin = proto.Float64(filter.RecoveryPriceMin.Value)
		}
		if filter.RecoveryPriceMax != nil {
			req.RecoveryPriceMax = proto.Float64(filter.RecoveryPriceMax.Value)
		}
		if filter.PriceRecoveryRatioMin != nil {
			req.PriceRecoveryRatioMin = proto.Float64(filter.PriceRecoveryRatioMin.Value)
		}
		if filter.PriceRecoveryRatioMax != nil {
			req.PriceRecoveryRatioMax = proto.Float64(filter.PriceRecoveryRatioMax.Value)
		}
	}
}
