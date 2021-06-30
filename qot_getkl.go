package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetkl"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotGetKL = 3006 //Qot_GetKL	获取 K 线
)

// 获取实时 K 线
func (api *FutuAPI) GetCurKL(ctx context.Context, security *Security, num int32, rehabType qotcommon.RehabType, klType qotcommon.KLType) (*RTKLine, error) {
	// 请求参数
	req := qotgetkl.Request{
		C2S: &qotgetkl.C2S{
			Security:  security.pb(),
			ReqNum:    &num,
			RehabType: (*int32)(&rehabType),
			KlType:    (*int32)(&klType),
		},
	}
	// 发送请求，同步返回结果
	ch := make(qotgetkl.ResponseChan)
	if err := api.get(ProtoIDQotGetKL, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return rtKLineFromGetPB(resp.GetS2C(), rehabType, klType), protocol.Error(resp)
	}
}

func rtKLineFromGetPB(pb *qotgetkl.S2C, rehabType qotcommon.RehabType, klType qotcommon.KLType) *RTKLine {
	if pb == nil {
		return nil
	}
	return &RTKLine{
		RehabType: rehabType,
		KLType:    klType,
		Security:  securityFromPB(pb.GetSecurity()),
		KLines:    kLineListFromPB(pb.GetKlList()),
	}
}
