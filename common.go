package futuapi

import (
	"github.com/hurisheng/go-futu-api/pb/common"
)

type ProgramStatus struct {
	Type       common.ProgramStatusType //*当前状态
	StrExtDesc string                   //额外描述
}

func programStatusFromPB(pb *common.ProgramStatus) *ProgramStatus {
	if pb == nil {
		return nil
	}
	return &ProgramStatus{
		Type:       pb.GetType(),
		StrExtDesc: pb.GetStrExtDesc(),
	}
}

type PacketID struct {
	ConnID   uint64 //当前 TCP 连接的连接 ID，一条连接的唯一标识，InitConnect 协议会返回
	SerialNo uint32 //自增序列号
}

func (p *PacketID) pb() *common.PacketID {
	return &common.PacketID{
		ConnID:   &p.ConnID,
		SerialNo: &p.SerialNo,
	}
}

func (api *FutuAPI) packetID() *PacketID {
	return &PacketID{
		ConnID:   api.ConnID(),
		SerialNo: api.serialNo(),
	}
}
