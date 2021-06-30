package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotrequesthistoryklquota"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotRequestHistoryKLQuota = 3104 //Qot_RequestHistoryKLQuota	获取历史 K 线额度
)

func (api *FutuAPI) GetHistoryKLQuota(ctx context.Context, detail bool) (*HistoryKLQuota, error) {
	req := qotrequesthistoryklquota.Request{
		C2S: &qotrequesthistoryklquota.C2S{
			BGetDetail: &detail,
		},
	}
	ch := make(qotrequesthistoryklquota.ResponseChan)
	if err := api.get(ProtoIDQotRequestHistoryKLQuota, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return historyKLQuotaFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type HistoryKLQuota struct {
	UsedQuota   int32                 //已使用过的额度，即当前周期内已经下载过多少只股票。
	RemainQuota int32                 //剩余额度
	DetailList  []*HistoryKLQuotaItem //每只拉取过的股票的下载时间
}

func historyKLQuotaFromPB(pb *qotrequesthistoryklquota.S2C) *HistoryKLQuota {
	if pb == nil {
		return nil
	}
	return &HistoryKLQuota{
		UsedQuota:   pb.GetUsedQuota(),
		RemainQuota: pb.GetRemainQuota(),
		DetailList:  historyKLQuotaItemListFromPB(pb.GetDetailList()),
	}
}

type HistoryKLQuotaItem struct {
	Security         *Security //拉取的股票
	RequestTime      string    //拉取的时间字符串
	RequestTimestamp int64     //拉取的时间戳
}

func historyKLQuotaItemFromPB(pb *qotrequesthistoryklquota.DetailItem) *HistoryKLQuotaItem {
	if pb == nil {
		return nil
	}
	return &HistoryKLQuotaItem{
		Security:         securityFromPB(pb.GetSecurity()),
		RequestTime:      pb.GetRequestTime(),
		RequestTimestamp: pb.GetRequestTimeStamp(),
	}
}

func historyKLQuotaItemListFromPB(pb []*qotrequesthistoryklquota.DetailItem) []*HistoryKLQuotaItem {
	if pb == nil {
		return nil
	}
	list := make([]*HistoryKLQuotaItem, len(pb))
	for i, v := range pb {
		list[i] = historyKLQuotaItemFromPB(v)
	}
	return list
}
