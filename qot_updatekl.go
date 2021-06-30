package futuapi

import (
	"context"

	"github.com/hurisheng/go-futu-api/pb/qotcommon"
	"github.com/hurisheng/go-futu-api/pb/qotupdatekl"
	"github.com/hurisheng/go-futu-api/protocol"
	"google.golang.org/protobuf/proto"
)

const (
	ProtoIDQotUpdateKL = 3007 //Qot_UpdateKL	推送 K 线
)

// 实时 K 线回调
func (api *FutuAPI) UpdateKL(ctx context.Context) (<-chan *UpdateKLResp, error) {
	ch := make(updateKLChan)
	if err := api.update(ProtoIDQotUpdateKL, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

type UpdateKLResp struct {
	KLine *RTKLine
	Err   error
}

type updateKLChan chan *UpdateKLResp

var _ protocol.RespChan = make(updateKLChan)

func (ch updateKLChan) Send(b []byte) error {
	var resp qotupdatekl.Response
	if err := proto.Unmarshal(b, &resp); err != nil {
		return err
	}
	ch <- &UpdateKLResp{
		KLine: rtKLineFromUpdatePB(resp.GetS2C()),
		Err:   protocol.Error(&resp),
	}
	return nil
}

func (ch updateKLChan) Close() {
	close(ch)
}

func rtKLineFromUpdatePB(pb *qotupdatekl.S2C) *RTKLine {
	if pb == nil {
		return nil
	}
	return &RTKLine{
		RehabType: qotcommon.RehabType(pb.GetRehabType()),
		KLType:    qotcommon.KLType(pb.GetKlType()),
		Security:  securityFromPB(pb.GetSecurity()),
		KLines:    kLineListFromPB(pb.GetKlList()),
	}
}
