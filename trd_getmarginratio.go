package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetmarginratio"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDTrdGetMarginRatio = 2223 // Trd_GetMarginRatio 获取融资融券数据

func init() {
	workers[ProtoIDTrdGetMarginRatio] = protocol.NewGetter()
}

// 获取融资融券数据
func (api *FutuAPI) GetMarginRatio(ctx context.Context, header *trdcommon.TrdHeader, securities []*qotcommon.Security) ([]*trdgetmarginratio.MarginRatioInfo, error) {

	if header == nil || len(securities) == 0 {
		return nil, ErrParameters
	}
	req := &trdgetmarginratio.Request{
		C2S: &trdgetmarginratio.C2S{
			Header:       header,
			SecurityList: securities,
		},
	}

	ch := make(chan *trdgetmarginratio.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetMarginRatio, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetMarginRatioInfoList(), protocol.Error(resp)
	}
}
