package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatert"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotUpdateRT = 3009 //Qot_UpdateRT	推送分时

func init() {
	workers[ProtoIDQotUpdateRT] = protocol.NewUpdater()
}

// 实时分时回调
func (api *FutuAPI) UpdateRT(ctx context.Context) (<-chan *qotupdatert.Response, error) {
	ch := make(chan *qotupdatert.Response)
	if err := api.proto.RegisterUpdate(ProtoIDQotUpdateRT, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	return ch, nil
}
