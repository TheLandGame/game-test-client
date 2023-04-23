package client

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

func (c *GameClient) SingIn() {
	reqMsg := &proto.Envelope{
		Type: proto.EnvelopeType_SigninPlayer,
		Payload: &proto.Envelope_SigninPlayerRequest{
			SigninPlayerRequest: &proto.SigninPlayerRequest{
				Token:      c.token,
				ClientTime: time_helper.NowUTCMill(),
				// SceneServiceAppId: "game-service-world-5",
			},
		},
	}
	c.net.Send(reqMsg)
}

func (c *GameClient) EnterMap() {
	reqMsg := &proto.Envelope{
		Type: proto.EnvelopeType_EnterMap,
		Payload: &proto.Envelope_EnterMapRequest{
			EnterMapRequest: &proto.EnterMapRequest{},
		},
	}
	c.net.Send(reqMsg)
}
