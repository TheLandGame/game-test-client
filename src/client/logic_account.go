package client

import (
	"fmt"
	"game-message-core/proto"
)



func (c *GameClient) QueryUser() {
	reqMsg := &proto.Envelope{
		Type: proto.EnvelopeType_QueryPlayer,
		Payload: &proto.Envelope_QueryPlayerRequest{
			QueryPlayerRequest: &proto.QueryPlayerRequest{
				Token: c.token,
			},
		},
	}
	c.net.Send(reqMsg)
}

func (c *GameClient) CreateUser() {
	reqMsg := &proto.Envelope{
		Type: proto.EnvelopeType_CreatePlayer,
		Payload: &proto.Envelope_CreatePlayerRequest{
			CreatePlayerRequest: &proto.CreatePlayerRequest{
				Token:    c.token,
				NickName: fmt.Sprintf("TEST_%d", c.userIdx),
				Icon:     "icon_avatar",
				RoleId:   1,
			},
		},
	}
	c.net.Send(reqMsg)
}
