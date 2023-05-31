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

	c.registerMsgHandler(proto.EnvelopeType_ItemGet, c.UserDataModel.LoadItemHandler)
	c.registerMsgHandler(proto.EnvelopeType_BroadCastInitItem, c.UserDataModel.InitItemHandler)

}

func (c *GameClient) InitMapElementHandler(packet *net_packet.NetPacket) {
	resp := &proto.BroadCastInitMapElementResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	// ids := make([]int64, len(resp.Entity))
	// for i, e := range resp.Entity {
	// 	ids[i] = e.TypeId.Id
	// }
	// serviceLog.Info("[%d] InitMapElement number %+v", c.userIdx, len(ids))
}

func (c *GameClient) PingHandler(packet *net_packet.NetPacket) {
	c.pingModel.OnResPing(packet)
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

		// 登录
		c.SingIn()
	}
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

	// 登录
	c.SingIn()
}

func (c *GameClient) SigninPlayerHandler(packet *net_packet.NetPacket) {
	resp := &proto.SigninPlayerResp{}
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

	c.UserDataModel.SetBaseData(resp.Player.BaseData)
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

	var spd float32 = 5.0
	for _, aid := range resp.Me.Profile {
		if client_ai.AttributeId(aid.Id) == client_ai.Attribute_MoveSpd {
			spd = float32(aid.Value) / 100.0 // 移动速度配置单位 CM 转换为 M
			break
		}
	}
	c.UserDataModel.SetSceneData(resp.Me)
	c.UserDataModel.MapId = resp.Location.MapId
	c.UserDataModel.Pos = &matrix.Vector3{
		X: float64(resp.Location.Loc.X),
		Y: float64(resp.Location.Loc.Y),
		Z: float64(resp.Location.Loc.Z),
	}
	c.UserDataModel.Dir = &matrix.Vector3{
		X: float64(resp.Me.Dir.X),
		Y: float64(resp.Me.Dir.Y),
		Z: float64(resp.Me.Dir.Z),
	}

	c.ClientAiModel.Init(
		&c.net, resp.Me.BaseData.UserId,
		resp.Location.MapId, c.UserDataModel.Pos,
		c.UserDataModel.Pos,
		spd,
	)
}

func (c *GameClient) UpdateSelfLocationHandler(packet *net_packet.NetPacket) {
	c.ClientAiModel.MoveModel.OnUpdateSelfLocationRes(packet)
}
