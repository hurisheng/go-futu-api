package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetipolist"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetIpoList = 3217 //Qot_GetIpoList	获取新股

func init() {
	workers[ProtoIDQotGetIpoList] = protocol.NewGetter()
}

// 获取 IPO 信息
func (api *FutuAPI) GetIPOList(ctx context.Context, market qotcommon.QotMarket) ([]*qotgetipolist.IpoData, error) {

	if market == qotcommon.QotMarket_QotMarket_Unknown {
		return nil, ErrParameters
	}
	req := &qotgetipolist.Request{
		C2S: &qotgetipolist.C2S{
			Market: proto.Int32(int32(market)),
		},
	}
	ch := make(chan *qotgetipolist.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetIpoList, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetIpoList(), protocol.Error(resp)
	}
}
