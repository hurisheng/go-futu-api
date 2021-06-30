package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotupdatert"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateRT = 3009 //Qot_UpdateRT	推送分时
)

// 实时分时回调
func (api *FutuAPI) UpdateRT(ctx context.Context) (<-chan *UpdateRTResp, error) {
	ch := make(updateRTChan)
	if err := api.update(ProtoIDQotUpdateRT, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateRTResp struct {
	RT  *RTData
	Err error
}

type updateRTChan chan *UpdateRTResp

var _ protocol.RespChan = make(updateRTChan)

func (ch updateRTChan) Send(b []byte) error {
	var resp qotupdatert.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateRTResp{
		RT:  rtDataFromUpdatePB(resp.GetS2C()),
		Err: protocol.Error(&resp),
	}
	return nil
}

func (ch updateRTChan) Close() {
	close(ch)
}

func rtDataFromUpdatePB(pb *qotupdatert.S2C) *RTData {
	if pb == nil {
		return nil
	}
	return &RTData{
		Security:   securityFromPB(pb.GetSecurity()),
		TimeShares: timeShareListFromPB(pb.GetRtList()),
	}
}
