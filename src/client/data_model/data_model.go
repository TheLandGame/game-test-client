package data_model

import (
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
)

const GET_DATA_CD_MS int64 = 1000

type UserDataModel struct {
	net *client_net.ClientNet

	baseData  *proto.PlayerBaseData
	sceneData *proto.Player

	MapId int32
	Pos   *matrix.Vector3
	Dir   *matrix.Vector3

	Items     []*proto.Item
	TalentExp []*proto.TalentExp

	preSendMsgMs int64
}

func (ai *UserDataModel) Init(net *client_net.ClientNet) {
	ai.net = net
}
func (p *UserDataModel) SetNet(net *client_net.ClientNet) {
	p.net = net
}

func (p *UserDataModel) SetBaseData(baseData *proto.PlayerBaseData) {
	p.baseData = baseData
}
func (p *UserDataModel) GetBaseData() *proto.PlayerBaseData {
	return p.baseData
}

func (p *UserDataModel) SetSceneData(data *proto.Player) {
	p.sceneData = data
}
func (p *UserDataModel) GetSceneData() *proto.Player {
	return p.sceneData
}

func (p *UserDataModel) LoadItem(curMs int64) {
	p.Items = []*proto.Item{}
	req := &proto.ItemGetReq{
		ReqTitle: &proto.ReqTitle{SeqId: p.net.NextSeqId()},
	}
	p.net.Send(proto.EnvelopeType_ItemGet, req)
	p.preSendMsgMs = curMs
}
func (p *UserDataModel) LoadItemHandler(packet *net_packet.NetPacket) {
	resp := &proto.ItemGetResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	p.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)

	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n",
			p.baseData.UserId, proto.EnvelopeType(packet.Id),
			resp.ResTitle.ErrorMessage,
		)
	}
}
func (p *UserDataModel) InitItemHandler(packet *net_packet.NetPacket) {
	resp := &proto.BroadCastInitItemResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n",
			p.baseData.UserId, proto.EnvelopeType(packet.Id),
			resp.ResTitle.ErrorMessage,
		)
	} else {
		p.Items = append(p.Items, resp.Items...)
	}
}

func (p *UserDataModel) LoadTalentExp(curMs int64) {
	req := &proto.QueryTalentExpReq{
		ReqTitle: &proto.ReqTitle{SeqId: p.net.NextSeqId()},
	}
	p.net.Send(proto.EnvelopeType_QueryTalentExp, req)
	p.preSendMsgMs = curMs
}
func (p *UserDataModel) LoadTalentExpHandler(packet *net_packet.NetPacket) {
	resp := &proto.QueryTalentExpResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	p.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)
	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n",
			p.baseData.UserId, proto.EnvelopeType(packet.Id),
			resp.ResTitle.ErrorMessage,
		)
	} else {
		p.TalentExp = resp.ExpData
	}
}

func (p *UserDataModel) LoadHomeAnimalList(curMs int64) {
	req := &proto.QueryAnimalListReq{
		ReqTitle: &proto.ReqTitle{SeqId: p.net.NextSeqId()},
	}
	p.net.Send(proto.EnvelopeType_QueryTalentExp, req)
	p.preSendMsgMs = curMs
}
func (p *UserDataModel) LoadHomeAnimalListHandler(packet *net_packet.NetPacket) {
	resp := &proto.QueryAnimalListResp{}
	err := protoTool.UnmarshalProto(packet.Body, resp)
	if err != nil {
		serviceLog.Error("Type:%v err: %v", proto.EnvelopeType(packet.Id), err.Error())
		return
	}

	p.net.PrintMsgUsedMs(proto.EnvelopeType(packet.Id), resp.ResTitle.SeqId)
	if resp.ResTitle.ErrorMessage != "" {
		serviceLog.Error("cli[%d] msg[%v] %s \n",
			p.baseData.UserId, proto.EnvelopeType(packet.Id),
			resp.ResTitle.ErrorMessage,
		)
	}
}

func (p *UserDataModel) Tick(curMs int64) {
	if p.preSendMsgMs+GET_DATA_CD_MS > curMs {
		return
	}
	p.LoadItem(curMs)
	p.LoadTalentExp(curMs)
	p.LoadHomeAnimalList(curMs)
}
