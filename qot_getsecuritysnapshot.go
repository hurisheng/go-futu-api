package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetsecuritysnapshot"
	"github.com/hurisheng/go-futu-api/protocol"
)

const ProtoIDQotGetSecuritySnapshot = 3203 //Qot_GetSecuritySnapshot	获取股票快照

func init() {
	workers[ProtoIDQotGetSecuritySnapshot] = protocol.NewGetter()
}

// 获取快照
func (api *FutuAPI) GetMarketSnapshot(ctx context.Context, securities []*qotcommon.Security) ([]*qotgetsecuritysnapshot.Snapshot, error) {
	if len(securities) == 0 {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotgetsecuritysnapshot.Request{
		C2S: &qotgetsecuritysnapshot.C2S{
			SecurityList: securities,
		},
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotgetsecuritysnapshot.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetSecuritySnapshot, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetSnapshotList(), protocol.Error(resp)
	}
}
