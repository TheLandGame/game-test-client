package client

import (
	"game-message-core/proto"
	"os"

	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

func (c *GameClient) SingIn() {
	sceneServiceAppId := os.Getenv("DEFAULT_SCENE_SER")

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
	if c.model != TEST_MODE_NORMAL {
		return
	}

	req := &proto.EnterMapReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
	}
	c.net.Send(proto.EnvelopeType_EnterMap, req)
	c.net.OnSendMsg(proto.EnvelopeType_EnterMap, req.ReqTitle.SeqId)
}
