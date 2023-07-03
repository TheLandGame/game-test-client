package client

import (
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
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
	c.registerMsgHandler(proto.EnvelopeType_Ping, c.pingModel.PingHandler)

	c.registerMsgHandler(proto.EnvelopeType_QueryPlayer, c.QueryPlayerHandler)
	c.registerMsgHandler(proto.EnvelopeType_CreatePlayer, c.CreatePlayerHandler)

	c.registerMsgHandler(proto.EnvelopeType_Login, c.LoginRespHandler)
	c.registerMsgHandler(proto.EnvelopeType_BroadCastQueue, c.BroadCastQueueRespHandler)
	c.registerMsgHandler(proto.EnvelopeType_SigninPlayer, c.SigninPlayerHandler)
	c.registerMsgHandler(proto.EnvelopeType_EnterMap, c.EnterMapHandler)
	c.registerMsgHandler(proto.EnvelopeType_UpdateSelfLocation, c.UpdateSelfLocationHandler)
	c.registerMsgHandler(proto.EnvelopeType_BroadCastInitMapElement, c.InitMapElementHandler)

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

func (c *GameClient) UpdateSelfLocationHandler(packet *net_packet.NetPacket) {
	c.ClientAiModel.MoveModel.OnUpdateSelfLocationRes(packet)
}
