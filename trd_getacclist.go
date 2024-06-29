package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdgetacclist"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdGetAccList = 2001 //Trd_GetAccList	获取业务账户列表

func init() {
	workers[ProtoIDTrdGetAccList] = protocol.NewGetter()
}

// 获取交易业务账户列表
func (api *FutuAPI) GetAccList(ctx context.Context,
	category trdcommon.TrdCategory, generalAcc *OptionalBool) ([]*trdcommon.TrdAcc, error) {

	// request information
	req := &trdgetacclist.Request{
		C2S: &trdgetacclist.C2S{
			UserID: proto.Uint64(api.UserID()),
		},
	}
	// optional parameters
	if category != trdcommon.TrdCategory_TrdCategory_Unknown {
		req.C2S.TrdCategory = proto.Int32(int32(category))
	}
	if generalAcc != nil {
		req.C2S.NeedGeneralSecAccount = proto.Bool(generalAcc.Value)
	}

	ch := make(chan *trdgetacclist.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdGetAccList, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetAccList(), protocol.Error(resp)
	}
}
