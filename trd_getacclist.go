package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/trdgetacclist"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDTrdGetAccList = 2001 //Trd_GetAccList	获取业务账户列表
)

// 获取交易业务账户列表
func (api *FutuAPI) GetAccList(ctx context.Context) ([]*TrdAcc, error) {
	req := trdgetacclist.Request{
		C2S: &trdgetacclist.C2S{
			UserID: &api.userID,
		},
	}
	ch := make(trdgetacclist.ResponseChan)
	if err := api.get(ProtoIDTrdGetAccList, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return trdAccListFromPB(resp.GetS2C().GetAccList()), protocol.Error(resp)
	}
}
