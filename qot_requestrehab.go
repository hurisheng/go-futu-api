package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotrequestrehab"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotRequestRehab = 3105 //Qot_RequestRehab	在线获取单只股票复权信息
)

// 获取复权因子
func (api *FutuAPI) GetRehab(ctx context.Context, security *Security) ([]*Rehab, error) {
	// 请求参数
	req := qotrequestrehab.Request{
		C2S: &qotrequestrehab.C2S{
			Security: security.pb(),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotrequestrehab.ResponseChan)
	if err := api.get(ProtoIDQotRequestRehab, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return rehabListFromPB(resp.GetS2C().GetRehabList()), protocol.Error(resp)
	}
}
