package client

import (
	"fmt"
	"game-message-core/proto"
)

func (c *GameClient) QueryUser() {
	req := &proto.QueryPlayerReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
		Token:    c.token,
	}
	c.net.Send(proto.EnvelopeType_QueryPlayer, req)
	c.net.OnSendMsg(proto.EnvelopeType_QueryPlayer, req.ReqTitle.SeqId)
}

func (c *GameClient) CreateUser() {
	req := &proto.CreatePlayerReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
		Token:    c.token,
		NickName: fmt.Sprintf("TEST_%d", c.userIdx),
		RoleId:   1,
		Gender:   "man",
		Icon:     "icon_avatar",
		Feature: &proto.PlayerFeature{
			Eyebrow: 123,
			Eye:     133,
			Hair:    162,
			Pants:   145,
			Skin:    109,
		},
	}
	c.net.Send(proto.EnvelopeType_CreatePlayer, req)
	c.net.OnSendMsg(proto.EnvelopeType_CreatePlayer, req.ReqTitle.SeqId)
}
