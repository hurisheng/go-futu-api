package futuapi

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/hurisheng/go-futu-api/pb/trdcommon"
	"github.com/hurisheng/go-futu-api/pb/trdunlocktrade"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDTrdUnlockTrade = 2005 //Trd_UnlockTrade	解锁或锁定交易

func init() {
	workers[ProtoIDTrdUnlockTrade] = protocol.NewGetter()
}

// 解锁交易
func (api *FutuAPI) UnlockTrade(ctx context.Context, unlock bool,
	pwd string, isMD5 bool, firm trdcommon.SecurityFirm) error {

	if unlock && pwd == "" {
		return ErrParameters
	}
	req := &trdunlocktrade.Request{
		C2S: &trdunlocktrade.C2S{
			Unlock: proto.Bool(unlock),
		},
	}
	if pwd != "" {
		if isMD5 {
			req.C2S.PwdMD5 = proto.String(pwd)
		} else {
			req.C2S.PwdMD5 = proto.String(fmt.Sprintf("%x", md5.Sum([]byte(pwd))))
		}
	}
	if firm != trdcommon.SecurityFirm_SecurityFirm_Unknown {
		req.C2S.SecurityFirm = proto.Int32(int32(firm))
	}

	ch := make(chan *trdunlocktrade.Response)
	if err := api.proto.RegisterGet(ProtoIDTrdUnlockTrade, req, protocol.NewProtobufChan(ch)); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return ErrChannelClosed
		}
		return protocol.Error(resp)
	}
}
