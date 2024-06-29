package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotgetreference"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotGetReference = 3206 //Qot_GetReference	获取正股相关股票

func init() {
	workers[ProtoIDQotGetReference] = protocol.NewGetter()
}

// 获取证券关联数据
func (api *FutuAPI) GetReferenceStockList(ctx context.Context, security *qotcommon.Security, refType qotgetreference.ReferenceType) ([]*qotcommon.SecurityStaticInfo, error) {

	if security == nil || refType == qotgetreference.ReferenceType_ReferenceType_Unknow {
		return nil, ErrParameters
	}
	req := &qotgetreference.Request{
		C2S: &qotgetreference.C2S{
			Security:      security,
			ReferenceType: proto.Int32(int32(refType)),
		},
	}
	ch := make(chan *qotgetreference.Response)
	if err := api.proto.RegisterGet(ProtoIDQotGetReference, req, protocol.NewProtobufChan(ch)); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return resp.GetS2C().GetStaticInfoList(), protocol.Error(resp)
	}
}
