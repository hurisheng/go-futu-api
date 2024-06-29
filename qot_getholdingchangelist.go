package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetholdingchangelist"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetHoldingChangeList = 3208 //Qot_GetHoldingChangeList 获取高管持股变动

func init() {
	workers[ProtoIDQotGetHoldingChangeList] = protocol.NewGetter()
}

// 获取高管持股变动
func (api *FutuAPI) GetHoldingChangeList(ctx context.Context, security *qotcommon.Security, holder qotcommon.HolderCategory,
	begin string, end string) (*qotgetholdingchangelist.S2C, error) {

	if security == nil || holder == qotcommon.HolderCategory_HolderCategory_Unknow {
		return nil, ErrParameters
	}
	req := &qotgetholdingchangelist.Request{
		C2S: &qotgetholdingchangelist.C2S{
			Security:       security,
			HolderCategory: proto.Int32(int32(holder)),
		},
	}
	if begin != "" {
		req.C2S.BeginTime = proto.String(begin)
	}
	if end != "" {
		req.C2S.EndTime = proto.String(end)
	}

	ch := make(chan *qotgetholdingchangelist.Response)
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
		return resp.GetS2C(), protocol.Error(resp)
	}
}
