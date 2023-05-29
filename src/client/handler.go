package client

import (
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/meland-client/src/client/client_ai"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

func (c *GameClient) onReceiveMsg(packet *net_packet.NetPacket) {
	eType := proto.EnvelopeType(packet.Id)

	if f, exist := c.msgEvent[eType]; exist {
		f(packet)
	} else {
		// serviceLog.Debug("cli[%d] Msg[%v] SeqId[%d] not register", c.userIdx, msg.Type, msg.SeqId)
	}
	net_packet.GetPool().Put(packet)
}

func (c *GameClient) registerMsgHandler(
	msgType proto.EnvelopeType, handler func(*net_packet.NetPacket),
) {
	c.msgEvent[msgType] = handler
}

func (c *GameClient) InitMsgHandler() {
	c.registerMsgHandler(proto.EnvelopeType_Ping, c.PingHandler)

	c.registerMsgHandler(proto.EnvelopeType_QueryPlayer, c.QueryPlayerHandler)
	c.registerMsgHandler(proto.EnvelopeType_CreatePlayer, c.CreatePlayerHandler)

	c.registerMsgHandler(proto.EnvelopeType_SigninPlayer, c.SigninPlayerHandler)
	c.registerMsgHandler(proto.EnvelopeType_EnterMap, c.EnterMapHandler)
	c.registerMsgHandler(proto.EnvelopeType_UpdateSelfLocation, c.UpdateSelfLocationHandler)
	c.registerMsgHandler(proto.EnvelopeType_BroadCastInitMapElement, c.InitMapElementHandler)

}

func (c *GameClient) InitMapElementHandler(packet *net_packet.NetPacket) {
	resp := &proto.BroadCastInitMapElementResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	ids := make([]int64, len(resp.Entity))
	for i, e := range resp.Entity {
		ids[i] = e.TypeId.Id
	}
	serviceLog.Info("InitMapElement %+v", ids)
}

func (c *GameClient) PingHandler(packet *net_packet.NetPacket) {
	c.pingModel.OnResPing(packet)
}

func (c *GameClient) QueryPlayerHandler(packet *net_packet.NetPacket) {
	resp := &proto.QueryPlayerResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)

	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error(
			"cli[%d] msg[%v] %s \n",
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
		c.SingIn()
	}
}

func (c *GameClient) CreatePlayerHandler(packet *net_packet.NetPacket) {
	resp := &proto.CreatePlayerResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)

	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error(
			"cli[%d] msg[%v] %s \n",
			c.userIdx, proto.EnvelopeType(packet.Id), resp.ResTitle.ErrorMessage,
		)
		c.stop()
		return
	}

	c.playerData.baseData = resp.Player
	// 登录
	c.SingIn()
}

func (c *GameClient) SigninPlayerHandler(packet *net_packet.NetPacket) {
	resp := &proto.SigninPlayerResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)

	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error(
			"cli[%d] msg[%v] %s \n",
			c.userIdx, proto.EnvelopeType(packet.Id), resp.ResTitle.ErrorMessage,
		)
		c.stop()
		return
	}

	c.playerData.sceneData = resp.Player
	time_helper.SetTimeOffsetMs(resp.ServerTime - resp.ClientTime)

	if resp.SceneServiceAppId == "" {
		serviceLog.Error("cli[%d] 无效 scene appId  \n", c.userIdx)
		c.stop()
	}
	c.EnterMap()
}

func (c *GameClient) EnterMapHandler(packet *net_packet.NetPacket) {
	resp := &proto.EnterMapResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	c.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)
	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error(
			"cli[%d] msg[%v] %s \n",
			c.userIdx, proto.EnvelopeType(packet.Id), resp.ResTitle.ErrorMessage,
		)
		c.stop()
		return
	}

	c.playerData.sceneData = resp.Me
	c.playerData.mapId = resp.Location.MapId
	c.playerData.Pos = &matrix.Vector3{
		X: float64(resp.Location.Loc.X),
		Y: float64(resp.Location.Loc.Y),
		Z: float64(resp.Location.Loc.Z),
	}
	c.playerData.Dir = &matrix.Vector3{
		X: float64(resp.Me.Dir.X),
		Y: float64(resp.Me.Dir.Y),
		Z: float64(resp.Me.Dir.Z),
	}

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

func (c *GameClient) UpdateSelfLocationHandler(packet *net_packet.NetPacket) {
	c.ClientAiModel.MoveModel.OnUpdateSelfLocationRes(packet)
}
