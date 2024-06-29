package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetstaticinfo"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetStaticInfo = 3202 //Qot_GetStaticInfo	获取股票静态信息

func init() {
	workers[ProtoIDQotGetStaticInfo] = protocol.NewGetter()
}

// 获取静态数据
func (api *FutuAPI) GetStockBasicInfo(ctx context.Context, market qotcommon.QotMarket, secType qotcommon.SecurityType, securities []*qotcommon.Security) ([]*qotcommon.SecurityStaticInfo, error) {
	if market == qotcommon.QotMarket_QotMarket_Unknown && len(securities) == 0 {
		return nil, ErrParameters
	}
	req := &qotgetstaticinfo.Request{
		C2S: &qotgetstaticinfo.C2S{
			SecType:      proto.Int32(int32(secType)),
			SecurityList: securities,
		},
	}
	if market != qotcommon.QotMarket_QotMarket_Unknown {
		req.C2S.Market = proto.Int32(int32(market))
	}
	ch := make(chan *qotgetstaticinfo.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetStaticInfo, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetStaticInfoList(), protocol.Error(resp)
	}
}
