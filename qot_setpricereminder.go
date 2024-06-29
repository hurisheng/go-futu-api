package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotsetpricereminder"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotSetPriceReminder = 3220 //Qot_SetPriceReminder	设置到价提醒

func init() {
	workers[ProtoIDQotSetPriceReminder] = protocol.NewGetter()
}

// 设置到价提醒
func (api *FutuAPI) SetPriceReminder(ctx context.Context, security *qotcommon.Security, op qotsetpricereminder.SetPriceReminderOp,
	key int64, remindType qotcommon.PriceReminderType, freq qotcommon.PriceReminderFreq, value *OptionalDouble, note string) (int64, error) {

	if security == nil || op == qotsetpricereminder.SetPriceReminderOp_SetPriceReminderOp_Unknown {
		return 0, ErrParameters
	}
	req := &qotsetpricereminder.Request{
		C2S: &qotsetpricereminder.C2S{
			Security: security,
			Op:       proto.Int32(int32(op)),
			Key:      proto.Int64(key),
		},
	}
	if remindType != qotcommon.PriceReminderType_PriceReminderType_Unknown {
		req.C2S.Type = proto.Int32(int32(remindType))
	}
	if freq != qotcommon.PriceReminderFreq_PriceReminderFreq_Unknown {
		req.C2S.Freq = proto.Int32(int32(freq))
	}
	if value != nil {
		req.C2S.Value = proto.Float64(value.Value)
	}
	if note != "" {
		req.C2S.Note = proto.String(note)
	}

	ch := make(chan *qotsetpricereminder.Response)
	if err := api.proto.RegisterGet(ProtoIDQotSetPriceReminder, req, protocol.NewProtobufChan(ch)); err != nil {
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
