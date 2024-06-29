package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetpositionlist"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdGetPositionList = 2102 //Trd_GetPositionList	获取账户持仓

func init() {
	workers[ProtoIDTrdGetPositionList] = protocol.NewGetter()
}

// 查询持仓
func (api *FutuAPI) GetPositionList(ctx context.Context, header *trdcommon.TrdHeader,
	filter *trdcommon.TrdFilterConditions, minPLRatio *OptionalDouble, maxPLRation *OptionalDouble, refresh *OptionalBool) ([]*trdcommon.Position, error) {

	if header == nil {
		return nil, ErrParameters
	}
	req := &trdgetpositionlist.Request{
		C2S: &trdgetpositionlist.C2S{
			Header:           header,
			FilterConditions: filter,
		},
	}
	if minPLRatio != nil {
		req.C2S.FilterPLRatioMin = proto.Float64(minPLRatio.Value)
	}
	if maxPLRation != nil {
		req.C2S.FilterPLRatioMax = proto.Float64(maxPLRation.Value)
	}
	if refresh != nil {
		req.C2S.RefreshCache = proto.Bool(refresh.Value)
	}

	ch := make(chan *trdgetpositionlist.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetPositionList, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetPositionList(), protocol.Error(resp)
	}
}
