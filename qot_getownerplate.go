package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetownerplate"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetOwnerPlate = 3207 //Qot_GetOwnerPlate	获取股票所属板块

func init() {
	workers[ProtoIDQotGetOwnerPlate] = protocol.NewGetter()
}

// 获取股票所属板块
func (api *FutuAPI) GetOwnerPlate(ctx context.Context, securities []*qotcommon.Security) ([]*qotgetownerplate.SecurityOwnerPlate, error) {

	if len(securities) == 0 {
		return nil, ErrParameters
	}
	// 请求数据
	req := &qotgetownerplate.Request{
		C2S: &qotgetownerplate.C2S{
			SecurityList: securities,
		},
	}
	ch := make(chan *qotgetownerplate.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetOwnerPlate, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetOwnerPlateList(), protocol.Error(resp)
	}
}
