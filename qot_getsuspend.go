package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetsuspend"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetSuspend = 3201 //Qot_GetSuspend 获取股票停牌信息

func init() {
	workers[ProtoIDQotGetSuspend] = protocol.NewGetter()
}

func (api *FutuAPI) GetSuspend(ctx context.Context, securities []*qotcommon.Security, begin string, end string) ([]*qotgetsuspend.SecuritySuspend, error) {
	if len(securities) == 0 || begin == "" || end == "" {
		return nil, ErrParameters
	}
	req := &qotgetsuspend.Request{
		C2S: &qotgetsuspend.C2S{
			SecurityList: securities,
			BeginTime:    proto.String(begin),
			EndTime:      proto.String(end),
		},
	}

	ch := make(chan *qotgetsuspend.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetSuspend, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetSecuritySuspendList(), protocol.Error(resp)
	}
}
