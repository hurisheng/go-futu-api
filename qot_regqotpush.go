package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotregqotpush"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotRegQotPush = 3002 // Qot_RegQotPush 注册推送

func init() {
	workers[ProtoIDQotRegQotPush] = protocol.NewGetter()
}

func (api *FutuAPI) RegQotPush(ctx context.Context, securities []*qotcommon.Security, subTypeList []qotcommon.SubType, isReg bool,
	rehabTypeList []qotcommon.RehabType, isFirstPush bool) error {

	if len(securities) == 0 || len(subTypeList) == 0 {
		return ErrParameters
	}
	req := &qotregqotpush.Request{
		C2S: &qotregqotpush.C2S{
			SecurityList: securities,
			SubTypeList:  make([]int32, len(subTypeList)),
			IsRegOrUnReg: proto.Bool(isReg),
			IsFirstPush:  proto.Bool(isFirstPush),
		},
	}
	for i, v := range subTypeList {
		req.C2S.SubTypeList[i] = int32(v)
	}
	if len(rehabTypeList) != 0 {
		req.C2S.RehabTypeList = make([]int32, len(rehabTypeList))
		for i, v := range rehabTypeList {
			req.C2S.RehabTypeList[i] = int32(v)
		}
	}

	ch := make(chan *qotregqotpush.Response)
	if err := api.proto.RegisterGet(ProtoIDQotRegQotPush, req, protocol.NewProtobufChan(ch)); err != nil {
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
