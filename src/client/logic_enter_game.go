package client

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

func (c *GameClient) SingIn() {
	// sceneServiceAppId := "game-service-world-735"
	sceneServiceAppId := ""

	req := &proto.SigninPlayerReq{
		ReqTitle:          &proto.ReqTitle{SeqId: c.net.NextSeqId()},
		Token:             c.token,
		ClientTime:        time_helper.NowUTCMill(),
		Reconnect:         false,
		SceneServiceAppId: sceneServiceAppId,
	}
	c.net.Send(proto.EnvelopeType_SigninPlayer, req)
	c.net.OnSendMsg(proto.EnvelopeType_SigninPlayer, req.ReqTitle.SeqId)
}

func (c *GameClient) EnterMap() {
	req := &proto.EnterMapReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
	}
	c.net.Send(proto.EnvelopeType_EnterMap, req)
	c.net.OnSendMsg(proto.EnvelopeType_EnterMap, req.ReqTitle.SeqId)
}
