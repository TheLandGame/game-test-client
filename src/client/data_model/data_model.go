package data_model

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
)

type DataModelState int

const (
	DATA_MODEL_STATE_READY DataModelState = iota
	DATA_MODEL_STATE_RUNNING
)

const GET_DATA_CD_MS int64 = 300

type UserDataModel struct {
	net *client_net.ClientNet

	baseData  *proto.PlayerBaseData
	sceneData *proto.Player

	state DataModelState

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

func (p *UserDataModel) SetState(state DataModelState) {
	p.state = state
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

func (p *UserDataModel) Tick(curMs int64) {
	if p.state != DATA_MODEL_STATE_RUNNING ||
		p.preSendMsgMs+GET_DATA_CD_MS > curMs {
		return
	}

	p.LoadItem(curMs)
	p.LoadTalentExp(curMs)
	p.LoadHomeAnimalList(curMs)
	p.QueryGranary(curMs)
}
