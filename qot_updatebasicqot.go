package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatebasicqot"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotUpdateBasicQot = 3005 //Qot_UpdateBasicQot	推送股票基本报价

)

func init() {
	workers[ProtoIDQotUpdateBasicQot] = protocol.NewUpdater()
}

// 实时报价回调
func (api *FutuAPI) UpdateBasicQot(ctx context.Context) (<-chan *qotupdatebasicqot.Response, error) {
	ch := make(chan *qotupdatebasicqot.Response)
	if err := api.proto.RegisterUpdate(ProtoIDQotUpdateBasicQot, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	return ch, nil
}
