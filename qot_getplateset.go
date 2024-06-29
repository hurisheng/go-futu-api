package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetplateset"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetPlateSet = 3204 //Qot_GetPlateSet	获取板块集合下的板块

func init() {
	workers[ProtoIDQotGetPlateSet] = protocol.NewGetter()
}

// 获取板块列表
func (api *FutuAPI) GetPlateList(ctx context.Context, market qotcommon.QotMarket, plateClass qotcommon.PlateSetType) ([]*qotcommon.PlateInfo, error) {

	if market == qotcommon.QotMarket_QotMarket_Unknown {
		return nil, ErrParameters
	}
	req := &qotgetplateset.Request{
		C2S: &qotgetplateset.C2S{
			Market:       proto.Int32(int32(market)),
			PlateSetType: proto.Int32(int32(plateClass)),
		},
	}
	ch := make(chan *qotgetplateset.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetPlateSet, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetPlateInfoList(), protocol.Error(resp)
	}
}
