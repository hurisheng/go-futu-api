package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotsetpricereminder"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotSetPriceReminder = 3220 //Qot_SetPriceReminder	设置到价提醒
)

// 设置到价提醒
func (api *FutuAPI) SetPriceReminder(ctx context.Context, security *Security, op qotsetpricereminder.SetPriceReminderOp,
	key int64, remindType qotcommon.PriceReminderType, freq qotcommon.PriceReminderFreq, value float64, note string) (int64, error) {
	req := qotsetpricereminder.Request{
		C2S: &qotsetpricereminder.C2S{
			Security: security.pb(),
			Op:       (*int32)(&op),
			Key:      &key,
			Type:     (*int32)(&remindType),
			Freq:     (*int32)(&freq),
			Value:    &value,
			Note:     &note,
		},
	}
	ch := make(qotsetpricereminder.ResponseChan)
	if err := api.get(ProtoIDQotSetPriceReminder, &req, ch); err != nil {
		return 0, err
	}
	select {
	case <-ctx.Done():
		return 0, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return 0, ErrChannelClosed
		}
		return resp.GetS2C().GetKey(), protocol.Error(resp)
	}
}
