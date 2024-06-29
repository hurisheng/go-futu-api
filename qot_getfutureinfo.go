package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetfutureinfo"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetFutureInfo = 3218 //Qot_GetFutureInfo	获取期货合约资料

func init() {
	workers[ProtoIDQotGetFutureInfo] = protocol.NewGetter()
}

// 获取期货合约资料
func (api *FutuAPI) GetFutureInfo(ctx context.Context, securities []*qotcommon.Security) ([]*qotgetfutureinfo.FutureInfo, error) {

	if len(securities) == 0 {
		return nil, ErrParameters
	}
	req := &qotgetfutureinfo.Request{
		C2S: &qotgetfutureinfo.C2S{
			SecurityList: securities,
		},
	}

	ch := make(chan *qotgetfutureinfo.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetFutureInfo, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetFutureInfoList(), protocol.Error(resp)
	}
}
