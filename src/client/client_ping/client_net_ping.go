package client_ping

import (
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
)

const PING_CD_MS int64 = 3000

type ClientPing struct {
	prePingAtMs int64
	net         *client_net.ClientNet
}

func (c *ClientPing) Init(clientNet *client_net.ClientNet) {
	c.net = clientNet
}

func (c *ClientPing) Tick(curMs int64) {
	if c.prePingAtMs+PING_CD_MS >= curMs {
		return
	}

	req := &proto.PingReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
	}
	c.net.Send(proto.EnvelopeType_Ping, req)
	c.net.OnSendMsg(proto.EnvelopeType_Ping, req.ReqTitle.SeqId)
	c.prePingAtMs = curMs
}

func (c *ClientPing) OnResPing(packet *net_packet.NetPacket) {
	resp := &proto.PingResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)
}
