package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotsub"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const ProtoIDQotSub = 3001 //Qot_Sub	订阅或者反订阅

func init() {
	workers[ProtoIDQotSub] = protocol.NewGetter()
}

// 订阅注册需要的实时信息，指定股票和订阅的数据类型即可。
// 香港市场（含正股、窝轮、牛熊、期权、期货）订阅，需要 LV1 及以上的权限，BMP 权限下不支持订阅。
func (api *FutuAPI) Subscribe(ctx context.Context, securities []*qotcommon.Security, subTypes []qotcommon.SubType,
	regPush bool, isFirstPush bool, isSubOrderBookDetail bool, extTime bool) error {
	return api.qotSub(ctx, securities, subTypes, true,
		&OptionalBool{regPush}, nil, &OptionalBool{isFirstPush}, &OptionalBool{false}, &OptionalBool{isSubOrderBookDetail}, &OptionalBool{extTime})
}

// 取消订阅
func (api *FutuAPI) Unsubscribe(ctx context.Context, securities []*qotcommon.Security, subTypes []qotcommon.SubType) error {
	return api.qotSub(ctx, securities, subTypes, false, nil, nil, nil, nil, nil, nil)
}

// 取消所有订阅
func (api *FutuAPI) UnsubscribeAll(ctx context.Context) error {
	return api.qotSub(ctx, nil, nil, false, nil, nil, nil, &OptionalBool{true}, nil, nil)
}

func (api *FutuAPI) qotSub(ctx context.Context, securities []*qotcommon.Security, subTypes []qotcommon.SubType, isSub bool,
	isRegPush *OptionalBool, rehabTypes []qotcommon.RehabType, isFirstPush *OptionalBool, isUnsubAll *OptionalBool, isSubOrderBookDetail *OptionalBool, extTime *OptionalBool) error {

	if (isUnsubAll == nil || !isUnsubAll.Value) && (len(securities) == 0 || len(subTypes) == 0) {
		return ErrParameters
	}
	// 拼装参数
	req := &qotsub.Request{
		C2S: &qotsub.C2S{
			SecurityList: securities,
			SubTypeList:  make([]int32, len(subTypes)),
			IsSubOrUnSub: proto.Bool(isSub),
		},
	}
	for i, v := range subTypes {
		req.C2S.SubTypeList[i] = int32(v)
	}
	if isRegPush != nil {
		req.C2S.IsRegOrUnRegPush = proto.Bool(isRegPush.Value)
	}
	if len(rehabTypes) != 0 {
		req.C2S.RegPushRehabTypeList = make([]int32, len(rehabTypes))
		for i, v := range rehabTypes {
			req.C2S.RegPushRehabTypeList[i] = int32(v)
		}
	}
	if isFirstPush != nil {
		req.C2S.IsFirstPush = proto.Bool(isFirstPush.Value)
	}
	if isUnsubAll != nil {
		req.C2S.IsUnsubAll = proto.Bool(isUnsubAll.Value)
	}
	if isSubOrderBookDetail != nil {
		req.C2S.IsSubOrderBookDetail = proto.Bool(isSubOrderBookDetail.Value)
	}
	if extTime != nil {
		req.C2S.ExtendedTime = proto.Bool(extTime.Value)
	}

	// 发送请求，同步返回结果
	ch := make(chan *qotsub.Response)
	if err := api.proto.RegisterGet(ProtoIDQotSub, req, protocol.NewProtobufChan(ch)); err != nil {
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
