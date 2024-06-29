package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetcodechange"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetCodeChange = 3216 //Qot_GetCodeChange 代码变换

func init() {
	workers[ProtoIDQotGetCodeChange] = protocol.NewGetter()
}

func (api *FutuAPI) GetCodeChange(ctx context.Context, placeholder int32,
	securities []*qotcommon.Security, timeFilter []*qotgetcodechange.TimeFilter, typeList []qotgetcodechange.CodeChangeType) ([]*qotgetcodechange.CodeChangeInfo, error) {

	req := &qotgetcodechange.Request{
		C2S: &qotgetcodechange.C2S{
			PlaceHolder: proto.Int32(placeholder),
		},
	}
	if len(securities) != 0 {
		req.C2S.SecurityList = securities
	}
	if len(timeFilter) != 0 {
		req.C2S.TimeFilterList = timeFilter
	}
	if len(typeList) != 0 {
		req.C2S.TypeList = make([]int32, len(typeList))
		for i, v := range typeList {
			req.C2S.TypeList[i] = int32(v)
		}
	}

	ch := make(chan *qotgetcodechange.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetCodeChange, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetCodeChangeList(), protocol.Error(resp)
	}
}
