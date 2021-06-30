package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdupdateorder"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDTrdUpdateOrder = 2208 //Trd_UpdateOrder	推送订单状态变动通知
)

// 响应订单推送回调
func (api *FutuAPI) UpdateOrder(ctx context.Context) (<-chan *UpdateOrderResp, error) {
	ch := make(updateOrderChan)
	if err := api.update(ProtoIDTrdUpdateOrder, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateOrderResp struct {
	Header *TrdHeader //交易公共参数头
	Order  *Order     //订单结构
	Err    error
}

type updateOrderChan chan *UpdateOrderResp

var _ protocol.RespChan = make(updateOrderChan)

func (ch updateOrderChan) Send(b []byte) error {
	var resp trdupdateorder.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateOrderResp{
		Header: trdHeaderFromPB(resp.GetS2C().GetHeader()),
		Order:  orderFromPB(resp.GetS2C().GetOrder()),
		Err:    protocol.Error(&resp),
	}
	return nil
}

func (ch updateOrderChan) Close() {
	close(ch)
}
