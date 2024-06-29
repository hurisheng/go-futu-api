package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdupdateorderfill"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDTrdUpdateOrderFill = 2218 //Trd_UpdateOrderFill	推送成交通知

func init() {
	workers[ProtoIDTrdUpdateOrderFill] = protocol.NewUpdater()
}

// 响应成交推送回调
func (api *FutuAPI) UpdateDeal(ctx context.Context) (<-chan *trdupdateorderfill.Response, error) {
	ch := make(chan *trdupdateorderfill.Response)
	if err := api.proto.RegisterUpdate(ProtoIDTrdUpdateOrderFill, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	return ch, nil
}
