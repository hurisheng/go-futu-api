package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotrequesthistorykl"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotRequestHistoryKL = 3103 //Qot_RequestHistoryKL	在线获取单只股票一段历史 K 线

func init() {
	workers[ProtoIDQotRequestHistoryKL] = protocol.NewGetter()
}

// 获取历史 K 线
func (api *FutuAPI) RequestHistoryKLine(ctx context.Context, security *qotcommon.Security, begin string, end string, klType qotcommon.KLType, rehabType qotcommon.RehabType,
	maxNum *OptionalInt32, klFields qotcommon.KLFields, nextKey []byte, extTime *OptionalBool) (*qotrequesthistorykl.S2C, error) {

	if security == nil || begin == "" || end == "" ||
		klType == qotcommon.KLType_KLType_Unknown || rehabType == qotcommon.RehabType_RehabType_None {
		return nil, ErrParameters
	}
	// 请求参数
	req := &qotrequesthistorykl.Request{
		C2S: &qotrequesthistorykl.C2S{
			RehabType:  proto.Int32(int32(rehabType)),
			KlType:     proto.Int32(int32(klType)),
			Security:   security,
			BeginTime:  proto.String(begin),
			EndTime:    proto.String(end),
			NextReqKey: nextKey,
		},
	}
	if maxNum != nil {
		req.C2S.MaxAckKLNum = proto.Int32(maxNum.Value)
	}
	if klFields != qotcommon.KLFields_KLFields_None {
		req.C2S.NeedKLFieldsFlag = proto.Int64(int64(klFields))
	}
	if extTime != nil {
		req.C2S.ExtendedTime = proto.Bool(extTime.Value)
	}
	// 发送请求，同步返回结果
	ch := make(chan *qotrequesthistorykl.Response)
	if err := api.proto.RegisterGet(ProtoIDQotRequestHistoryKL, req, protocol.NewProtobufChan(ch)); err != nil {
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
