package client

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/meland-client/src/client/client_ai"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

func (c *GameClient) onReceiveMsg(msg *proto.Envelope) {
	if f, exist := c.msgEvent[msg.Type]; exist {
		f(msg)
	} else {
		// serviceLog.Debug("cli[%d] Msg[%v] SeqId[%d] not register", c.userIdx, msg.Type, msg.SeqId)
	}
}

func (c *GameClient) InitMsgHandler() {
	c.registerMsgHandler(proto.EnvelopeType_BroadCastMsgAggregation, c.AggregationMsgHandler)

	c.registerMsgHandler(proto.EnvelopeType_Ping, c.PingHandler)

	c.registerMsgHandler(proto.EnvelopeType_QueryPlayer, c.QueryPlayerHandler)
	c.registerMsgHandler(proto.EnvelopeType_CreatePlayer, c.CreatePlayerHandler)

	c.registerMsgHandler(proto.EnvelopeType_SigninPlayer, c.SigninPlayerHandler)
	c.registerMsgHandler(proto.EnvelopeType_EnterMap, c.EnterMapHandler)
	c.registerMsgHandler(proto.EnvelopeType_UpdateSelfLocation, c.UpdateSelfLocationHandler)

}

func (c *GameClient) AggregationMsgHandler(msg *proto.Envelope) {
	res := msg.GetBroadCastMsgAggregationResponse()
	for _, aMsg := range res.MessageList {
		c.onReceiveMsg(aMsg)
	}
}

func (c *GameClient) PingHandler(msg *proto.Envelope) {
	c.pingModel.OnResPing(msg)
}

func (c *GameClient) QueryPlayerHandler(msg *proto.Envelope) {
	if msg.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n", c.userIdx, msg.Type, msg.ErrorMessage)
		return
	}

	res := msg.GetQueryPlayerResponse()
	if res.Player == nil || res.Player.UserId <= 0 {
		c.CreateUser() // create player

	} else {
		c.playerData.baseData = res.Player
		// 登录
		c.SingIn()
	}
}

func (c *GameClient) CreatePlayerHandler(msg *proto.Envelope) {
	if msg.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n", c.userIdx, msg.Type, msg.ErrorMessage)
		return
	}

	res := msg.GetCreatePlayerResponse()
	c.playerData.baseData = res.Player

	// 登录
	c.SingIn()
}

func (c *GameClient) SigninPlayerHandler(msg *proto.Envelope) {
	if msg.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n", c.userIdx, msg.Type, msg.ErrorMessage)
		c.stop()
		return
	}

	res := msg.GetSigninPlayerResponse()
	c.playerData.sceneData = res.Player
	time_helper.SetTimeOffsetMs(res.ServerTime - res.ClientTime)

	if res.SceneServiceAppId == "" {
		serviceLog.Error("cli[%d] 无效 scene appId  \n", c.userIdx)
		c.stop()
	}

	c.EnterMap()

}

func (c *GameClient) EnterMapHandler(msg *proto.Envelope) {
	if msg.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n", c.userIdx, msg.Type, msg.ErrorMessage)
		c.stop()
	}

	res := msg.GetEnterMapResponse()
	c.playerData.sceneData = res.Me
	c.playerData.Pos = &matrix.Vector3{X: float64(res.Location.Loc.X), Y: float64(res.Location.Loc.Y), Z: float64(res.Location.Loc.Z)}
	c.playerData.mapId = res.Location.MapId
	c.playerData.Dir = &matrix.Vector3{X: float64(res.Me.Dir.X), Y: float64(res.Me.Dir.Y), Z: float64(res.Me.Dir.Z)}

	var spd float32 = 5.0
	for _, aid := range c.playerData.sceneData.Profile {
		if client_ai.AttributeId(aid.Id) == client_ai.Attribute_MoveSpd {
			spd = float32(aid.Value) / 100.0 // 移动速度配置单位 CM 转换为 M
			break
		}
	}
	c.ClientAiModel.Init(
		&c.net, c.playerData.baseData.UserId,
		c.playerData.mapId, c.playerData.Pos, c.playerData.Dir, spd,
	)
}

func (c *GameClient) UpdateSelfLocationHandler(msg *proto.Envelope) {
	if msg.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n", c.userIdx, msg.Type, msg.ErrorMessage)
		return
	}

	c.ClientAiModel.MoveModel.OnUpdateSelfLocationRes(msg)
}
