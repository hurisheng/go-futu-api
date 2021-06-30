package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotrequesthistorykl"
	"github.com/hurisheng/go-futu-api/protocol"
)

const (
	ProtoIDQotRequestHistoryKL = 3103 //Qot_RequestHistoryKL	在线获取单只股票一段历史 K 线
)

// 获取历史 K 线
func (api *FutuAPI) RequestHistoryKLine(ctx context.Context, security *Security, begin string, end string, klType qotcommon.KLType, rehabType qotcommon.RehabType,
	maxNum int32, fields qotcommon.KLFields, nextKey []byte, extTime bool) (*HistoryKLine, error) {
	// 请求参数
	req := qotrequesthistorykl.Request{
		C2S: &qotrequesthistorykl.C2S{
			RehabType:    (*int32)(&rehabType),
			KlType:       (*int32)(&klType),
			Security:     security.pb(),
			BeginTime:    &begin,
			EndTime:      &end,
			NextReqKey:   nextKey,
			ExtendedTime: &extTime,
		},
	}
	if maxNum != 0 {
		req.C2S.MaxAckKLNum = &maxNum
	}
	if fields != 0 {
		var klFields int64 = int64(fields)
		req.C2S.NeedKLFieldsFlag = &klFields
	}
	// 发送请求，同步返回结果
	ch := make(qotrequesthistorykl.ResponseChan)
	if err := api.get(ProtoIDQotRequestHistoryKL, &req, ch); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-ch:
		if !ok {
			return nil, ErrChannelClosed
		}
		return historyKLineFromPB(resp.GetS2C()), protocol.Error(resp)
	}
}

type HistoryKLine struct {
	Security *Security //证券
	KLines   []*KLine  //K 线数据
	NextKey  []byte    //分页请求 key。一次请求没有返回所有数据时，下次请求带上这个 key，会接着请求
}

func historyKLineFromPB(pb *qotrequesthistorykl.S2C) *HistoryKLine {
	if pb == nil {
		return nil
	}
	return &HistoryKLine{
		Security: securityFromPB(pb.GetSecurity()),
		KLines:   kLineListFromPB(pb.GetKlList()),
		NextKey:  nextKeyFromPB(pb.GetNextReqKey()),
	}
}

func nextKeyFromPB(pb []byte) []byte {
	if pb == nil {
		return nil
	}
	k := make([]byte, len(pb))
	copy(k, pb)
	return k
}
