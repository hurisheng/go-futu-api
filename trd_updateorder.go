package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdupdateorder"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDTrdUpdateOrder = 2208 //Trd_UpdateOrder	推送订单状态变动通知

func init() {
	workers[ProtoIDTrdUpdateOrder] = protocol.NewUpdater()
}

// 响应订单推送回调
func (api *FutuAPI) UpdateOrder(ctx context.Context) (<-chan *trdupdateorder.Response, error) {
	ch := make(chan *trdupdateorder.Response)
	if err := api.proto.RegisterUpdate(ProtoIDTrdUpdateOrder, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	return ch, nil
}
