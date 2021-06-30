package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotgetcapitaldistribution"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetCapitalDistribution = 3212 //Qot_GetCapitalDistribution
)

// 获取资金分布
func (api *FutuAPI) GetCapitalDistribution(ctx context.Context, security *Security) (*CapitalDistribution, error) {
	// 请求数据
	req := qotgetcapitaldistribution.Request{
		C2S: &qotgetcapitaldistribution.C2S{
			Security: security.pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetcapitaldistribution.ResponseChan)
	if err := api.get(ProtoIDQotGetCapitalDistribution, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return capitalDistributionFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

// 根据历史成交数据将逐笔成交记录划分成大单，中单，小单。以正股前一个月（或窝轮前三天）的平均每笔成交额为参考值，小于该平均值为小单，大于等于该金额的10倍为大单，其余为中单。
type CapitalDistribution struct {
	InBig           float64 //流入资金额度，大单
	InMid           float64 //流入资金额度，中单
	InSmall         float64 //流入资金额度，小单
	OutBig          float64 //流出资金额度，大单
	OutMid          float64 //流出资金额度，中单
	OutSmall        float64 //流出资金额度，小单
	UpdateTime      string  //更新时间字符串
	UpdateTimestamp float64 //更新时间戳
}

func capitalDistributionFromPB(pb *qotgetcapitaldistribution.S2C) *CapitalDistribution {
	if pb == nil {
		return nil
	}
	return &CapitalDistribution{
		InBig:           pb.GetCapitalInBig(),
		InMid:           pb.GetCapitalInMid(),
		InSmall:         pb.GetCapitalInSmall(),
		OutBig:          pb.GetCapitalOutBig(),
		OutMid:          pb.GetCapitalOutMid(),
		OutSmall:        pb.GetCapitalOutSmall(),
		UpdateTime:      pb.GetUpdateTime(),
		UpdateTimestamp: pb.GetUpdateTimestamp(),
	}
}
