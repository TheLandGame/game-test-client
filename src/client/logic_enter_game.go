package client

import (
	"game-message-core/proto"
	"game-message-core/protoTool"
	"os"

	"github.com/Meland-Inc/meland-client/src/client/client_ai"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

func (c *GameClient) Login() {
	req := &proto.LoginReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
		Token:    c.token,
	}
	c.net.Send(proto.EnvelopeType_Login, req)
	c.net.OnSendMsg(proto.EnvelopeType_Login, req.ReqTitle.SeqId)
}

func (c *GameClient) LoginRespHandler(packet *net_packet.NetPacket) {
	resp := &proto.LoginResp{}
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
	serviceLog.Info("%d login response queue = %d", c.userIdx, resp.Queue)
}

func (c *GameClient) BroadCastQueueRespHandler(packet *net_packet.NetPacket) {
	resp := &proto.BroadCastQueueResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	serviceLog.Info(
		"receive BroadCastQueue user [%d]  state[%v], queue[%d]",
		c.userIdx, resp.LoginState, resp.Queue,
	)

	if resp.LoginState == proto.LoginState_SingIn {
		c.SingIn()
	}

}

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

	c.playerData.sceneData = resp.Player
	time_helper.SetTimeOffsetMs(resp.ServerTime - resp.ClientTime)

	if resp.SceneServiceAppId == "" {
		serviceLog.Error("cli[%d] 无效 scene appId  \n", c.userIdx)
		c.stop()
		return
	}

	c.EnterMap()
}

func (c *GameClient) EnterMap() {
	req := &proto.EnterMapReq{
		ReqTitle: &proto.ReqTitle{SeqId: c.net.NextSeqId()},
	}
	c.net.Send(proto.EnvelopeType_EnterMap, req)
	c.net.OnSendMsg(proto.EnvelopeType_EnterMap, req.ReqTitle.SeqId)
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
