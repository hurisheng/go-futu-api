package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetpricereminder"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetPriceReminder = 3221 //Qot_GetPriceReminder	获取到价提醒
)

// 获取到价提醒列表
func (api *FutuAPI) GetPriceReminder(ctx context.Context, security *Security, market qotcommon.QotMarket) ([]*PriceReminder, error) {
	req := qotgetpricereminder.Request{
		C2S: &qotgetpricereminder.C2S{
			Security: security.pb(),
			Market:   (*int32)(&market),
		},
	}
	ch := make(qotgetpricereminder.ResponseChan)
	if err := api.get(ProtoIDQotGetPriceReminder, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return priceReminderListFromPB(resp.GetS2C().GetPriceReminderList()), protocol.Error(resp)
	}
}

type PriceReminder struct {
	Security *Security            // 股票
	ItemList []*PriceReminderItem // 提醒信息列表
}

func priceReminderFromPB(pb *qotgetpricereminder.PriceReminder) *PriceReminder {
	if pb == nil {
		return nil
	}
	return &PriceReminder{
		Security: securityFromPB(pb.GetSecurity()),
		ItemList: priceReminderItemListFromPB(pb.GetItemList()),
	}
}

func priceReminderListFromPB(pb []*qotgetpricereminder.PriceReminder) []*PriceReminder {
	if pb == nil {
		return nil
	}
	list := make([]*PriceReminder, len(pb))
	for i, v := range pb {
		list[i] = priceReminderFromPB(v)
	}
	return list
}

type PriceReminderItem struct {
	Key      int64                       // 每个提醒的唯一标识
	ItemType qotcommon.PriceReminderType // Qot_Common::PriceReminderType 提醒类型
	Value    float64                     // 提醒参数值
	Note     string                      // 备注
	Freq     qotcommon.PriceReminderFreq // Qot_Common::PriceReminderFreq 提醒频率类型
	IsEnable bool                        // 该提醒设置是否生效。false 不生效，true 生效
}

func priceReminderItemFromPB(pb *qotgetpricereminder.PriceReminderItem) *PriceReminderItem {
	if pb == nil {
		return nil
	}
	return &PriceReminderItem{
		Key:      pb.GetKey(),
		ItemType: qotcommon.PriceReminderType(pb.GetType()),
		Value:    pb.GetValue(),
		Note:     pb.GetNote(),
		Freq:     qotcommon.PriceReminderFreq(pb.GetFreq()),
		IsEnable: pb.GetIsEnable(),
	}
}

func priceReminderItemListFromPB(pb []*qotgetpricereminder.PriceReminderItem) []*PriceReminderItem {
	if pb == nil {
		return nil
	}
	list := make([]*PriceReminderItem, len(pb))
	for i, v := range pb {
		list[i] = priceReminderItemFromPB(v)
	}
	return list
}
