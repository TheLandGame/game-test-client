package client_ping

import (
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
)

func (c *ClientPing) ping() {
	req := &proto.PingReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
	}
	c.net.Send(proto.EnvelopeType_Ping, req)
	c.net.OnSendMsg(proto.EnvelopeType_Ping, req.ReqTitle.SeqId)
}

func (c *ClientPing) PingHandler(packet *net_packet.NetPacket) {
	resp := &proto.PingResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}
	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)
}
