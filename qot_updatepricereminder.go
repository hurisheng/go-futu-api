package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotupdatepricereminder"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdatePriceReminder = 3019 //Qot_UpdatePriceReminder	到价提醒通知
)

// 到价提醒回调
func (api *FutuAPI) UpdatePriceReminder(ctx context.Context) (<-chan *UpdatePriceReminderResp, error) {
	ch := make(updatePriceReminderChan)
	if err := api.update(ProtoIDQotUpdatePriceReminder, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type RTPriceReminder struct {
	Security     *Security                           //股票
	Price        float64                             //价格
	ChangeRate   float64                             //当日涨跌幅
	MarketStatus qotupdatepricereminder.MarketStatus //市场状态
	Content      string                              //内容
	Note         string                              //备注
	Key          int64                               //到价提醒的标识
	Type         qotcommon.PriceReminderType         //Qot_Common::PriceReminderType，提醒频率类型
	SetValue     float64                             //设置的提醒值
	CurValue     float64                             //设置的提醒类型触发时当前值
}

func rtPriceReminderFromPB(pb *qotupdatepricereminder.S2C) *RTPriceReminder {
	if pb == nil {
		return nil
	}
	return &RTPriceReminder{
		Security:     securityFromPB(pb.GetSecurity()),
		Price:        pb.GetPrice(),
		ChangeRate:   pb.GetChangeRate(),
		MarketStatus: qotupdatepricereminder.MarketStatus(pb.GetMarketStatus()),
		Content:      pb.GetContent(),
		Note:         pb.GetNote(),
		Key:          pb.GetKey(),
		Type:         qotcommon.PriceReminderType(pb.GetType()),
		SetValue:     pb.GetSetValue(),
		CurValue:     pb.GetCurValue(),
	}
}

type UpdatePriceReminderResp struct {
	Reminder *RTPriceReminder
	Err      error
}

type updatePriceReminderChan chan *UpdatePriceReminderResp

var _ protocol.RespChan = make(updatePriceReminderChan)

func (ch updatePriceReminderChan) Close() {
	close(ch)
}

func (ch updatePriceReminderChan) Send(b []byte) error {
	var resp qotupdatepricereminder.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdatePriceReminderResp{
		Reminder: rtPriceReminderFromPB(resp.GetS2C()),
		Err:      protocol.Error(&resp),
	}
	return nil
}
