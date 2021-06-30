package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdgetmarginratio"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdGetMarginRatio = 2223 // Trd_GetMarginRatio 获取融资融券数据
)

// 获取融资融券数据
func (api *FutuAPI) GetMarginRatio(ctx context.Context, header *TrdHeader, securities []*Security) ([]*MarginRatio, error) {
	req := trdgetmarginratio.Request{
		C2S: &trdgetmarginratio.C2S{
			Header:       header.pb(),
			SecurityList: securityList(securities).pb(),
		},
	}
	ch := make(trdgetmarginratio.ResponseChan)
	if err := api.get(ProtoIDTrdGetMarginRatio, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return marginRatioListFromPB(resp.GetS2C().GetMarginRatioInfoList()), protocol.Error(resp)
	}
}

type MarginRatio struct {
	Security        *Security //股票
	IsLongPermit    bool      //是否允许融资
	IsShortPermit   bool      //是否允许融券
	ShortPoolRemain float64   //卖空池剩余（股）
	ShortFeeRate    float64   //融券参考利率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	AlertLongRatio  float64   //融资预警比率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	AlertShortRatio float64   //融券预警比率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	IMLongRatio     float64   //融资初始保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	IMShortRatio    float64   //融券初始保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	MCMLongRatio    float64   //融资 margin call 保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	MCMShortRatio   float64   //融券 margin call 保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	MMLongRatio     float64   //融资维持保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
	MMShortRatio    float64   //融券维持保证金率（该字段为百分比字段，默认不展示 %，如 20 实际对应 20%）
}

func marginRatioFromPB(pb *trdgetmarginratio.MarginRatioInfo) *MarginRatio {
	if pb == nil {
		return nil
	}
	return &MarginRatio{
		Security:        securityFromPB(pb.GetSecurity()),
		IsLongPermit:    pb.GetIsLongPermit(),
		IsShortPermit:   pb.GetIsShortPermit(),
		ShortPoolRemain: pb.GetShortPoolRemain(),
		ShortFeeRate:    pb.GetShortFeeRate(),
		AlertLongRatio:  pb.GetAlertLongRatio(),
		AlertShortRatio: pb.GetAlertShortRatio(),
		IMLongRatio:     pb.GetImLongRatio(),
		IMShortRatio:    pb.GetImShortRatio(),
		MCMLongRatio:    pb.GetMcmLongRatio(),
		MCMShortRatio:   pb.GetMcmShortRatio(),
		MMLongRatio:     pb.GetMmLongRatio(),
		MMShortRatio:    pb.GetMmShortRatio(),
	}
}

func marginRatioListFromPB(pb []*trdgetmarginratio.MarginRatioInfo) []*MarginRatio {
	if pb == nil {
		return nil
	}
	list := make([]*MarginRatio, len(pb))
	for i, v := range pb {
		list[i] = marginRatioFromPB(v)
	}
	return list
}
