package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetoptionchain"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetOptionChain = 3209 //Qot_GetOptionChain	获取期权链

func init() {
	workers[ProtoIDQotGetOptionChain] = protocol.NewGetter()
}

// 获取期权链
func (api *FutuAPI) GetOptionChain(ctx context.Context, owner *qotcommon.Security, begin string, end string,
	indexOptType qotcommon.IndexOptionType, optType qotcommon.OptionType, cond qotgetoptionchain.OptionCondType,
	filter *qotgetoptionchain.DataFilter) ([]*qotgetoptionchain.OptionChain, error) {
	if owner == nil || begin == "" || end == "" {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetoptionchain.Request{
		C2S: &qotgetoptionchain.C2S{
			Owner:      owner,
			BeginTime:  proto.String(begin),
			EndTime:    proto.String(end),
			DataFilter: filter,
		},
	}

	if indexOptType != qotcommon.IndexOptionType_IndexOptionType_Unknown {
		req.C2S.IndexOptionType = proto.Int32(int32(indexOptType))
	}
	if optType != qotcommon.OptionType_OptionType_Unknown {
		req.C2S.Type = proto.Int32(int32(optType))
	}
	if cond != qotgetoptionchain.OptionCondType_OptionCondType_Unknow {
		req.C2S.Condition = proto.Int32(int32(cond))
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetoptionchain.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetOptionChain, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetOptionChain(), protocol.Error(resp)
	}
}
