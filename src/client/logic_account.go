package client

import (
	"fmt"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
)

func (c *GameClient) QueryUser() {
	req := &proto.QueryPlayerReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
		Token:    c.token,
	}
	c.net.Send(proto.EnvelopeType_QueryPlayer, req)
	c.net.OnSendMsg(proto.EnvelopeType_QueryPlayer, req.ReqTitle.SeqId)
}

func (c *GameClient) QueryPlayerHandler(packet *net_packet.NetPacket) {
	resp := &proto.QueryPlayerResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)

	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n",
			c.userIdx, proto.EnvelopeType(packet.Id), resp.ResTitle.ErrorMessage,
		)
		c.stop()
		return
	}

	if resp.Player == nil || resp.Player.UserId <= 0 {
		c.CreateUser() // create player
	} else {
		c.playerData.baseData = resp.Player
		// 登录
		c.Login()
	}
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

func (c *GameClient) CreatePlayerHandler(packet *net_packet.NetPacket) {
	resp := &proto.CreatePlayerResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)

	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n",
			c.userIdx, proto.EnvelopeType(packet.Id), resp.ResTitle.ErrorMessage,
		)
		c.stop()
		return
	}

	c.playerData.baseData = resp.Player
	// 登录
	c.Login()
}
