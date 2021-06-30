package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatebasicqot"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateBasicQot = 3005 //Qot_UpdateBasicQot	推送股票基本报价
)

// 实时报价回调
func (api *FutuAPI) UpdateBasicQot(ctx context.Context) (<-chan *UpdateBasicQotResp, error) {
	ch := make(updateBasicQotChan)
	if err := api.update(ProtoIDQotUpdateBasicQot, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateBasicQotResp struct {
	BasicQot []*BasicQot //股票基本行情
	Err      error
}

type updateBasicQotChan chan *UpdateBasicQotResp

var _ protocol.RespChan = make(updateBasicQotChan)

func (ch updateBasicQotChan) Send(b []byte) error {
	var resp qotupdatebasicqot.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateBasicQotResp{
		BasicQot: basicQotListFromPB(resp.GetS2C().GetBasicQotList()),
		Err:      protocol.Error(&resp),
	}
	return nil
}

func (ch updateBasicQotChan) Close() {
	close(ch)
}
