package data_model

import (
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
)

func (p *UserDataModel) LoadItem(curMs int64) {
	p.Items = []*proto.Item{}
	req := &proto.ItemGetReq{
		ReqTitle: &proto.ReqTitle{SeqId: p.net.NextSeqId()},
	}
	p.net.Send(proto.EnvelopeType_ItemGet, req)
	p.net.OnSendMsg(proto.EnvelopeType_ItemGet, req.ReqTitle.SeqId)
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
	p.net.OnSendMsg(proto.EnvelopeType_QueryTalentExp, req.ReqTitle.SeqId)
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
	p.net.Send(proto.EnvelopeType_QueryAnimalList, req)
	p.net.OnSendMsg(proto.EnvelopeType_QueryAnimalList, req.ReqTitle.SeqId)
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

func (p *UserDataModel) QueryGranary(curMs int64) {
	req := &proto.QueryGranaryReq{
		ReqTitle: &proto.ReqTitle{SeqId: p.net.NextSeqId()},
	}
	p.net.Send(proto.EnvelopeType_QueryGranary, req)
	p.net.OnSendMsg(proto.EnvelopeType_QueryGranary, req.ReqTitle.SeqId)
	p.preSendMsgMs = curMs
}
func (p *UserDataModel) QueryGranaryHandler(packet *net_packet.NetPacket) {
	resp := &proto.QueryGranaryResp{}
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
