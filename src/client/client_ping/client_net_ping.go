package client_ping

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/net_util"
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

	msg := net_util.MakePingMsg()
	c.net.Send(msg)
	c.prePingAtMs = curMs
}

func (c *ClientPing) OnResPing(msg *proto.Envelope) {

}
