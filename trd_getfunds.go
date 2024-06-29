package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetfunds"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdGetFunds = 2101 //Trd_GetFunds	获取账户资金

func init() {
	workers[ProtoIDTrdGetFunds] = protocol.NewGetter()
}

// 查询账户资金
func (api *FutuAPI) GetFunds(ctx context.Context, header *trdcommon.TrdHeader,
	refresh *OptionalBool, currency trdcommon.Currency) (*trdcommon.Funds, error) {

	if header == nil {
		return nil, ErrParameters
	}
	req := &trdgetfunds.Request{
		C2S: &trdgetfunds.C2S{
			Header: header,
		},
	}
	if refresh != nil {
		req.C2S.RefreshCache = proto.Bool(refresh.Value)
	}
	if currency != trdcommon.Currency_Currency_Unknown {
		req.C2S.Currency = proto.Int32(int32(currency))
	}

	ch := make(chan *trdgetfunds.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetFunds, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetFunds(), protocol.Error(resp)
	}
}
