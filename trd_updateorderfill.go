package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdupdateorderfill"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDTrdUpdateOrderFill = 2218 //Trd_UpdateOrderFill	推送成交通知
)

// 响应成交推送回调
func (api *FutuAPI) UpdateDeal(ctx context.Context) (<-chan *UpdateDealResp, error) {
	ch := make(updateDealChan)
	if err := api.update(ProtoIDTrdUpdateOrderFill, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateDealResp struct {
	Header    *TrdHeader
	OrderFill *OrderFill
	Err       error
}

type updateDealChan chan *UpdateDealResp

var _ protocol.RespChan = make(updateDealChan)

func (ch updateDealChan) Send(b []byte) error {
	var resp trdupdateorderfill.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateDealResp{
		Header:    trdHeaderFromPB(resp.GetS2C().GetHeader()),
		OrderFill: orderFillFromPB(resp.GetS2C().GetOrderFill()),
		Err:       protocol.Error(&resp),
	}
	return nil
}

func (ch updateDealChan) Close() {
	close(ch)
}
