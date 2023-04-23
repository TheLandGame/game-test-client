package move_model

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
	"github.com/Meland-Inc/meland-client/src/common/random"
)

var birthPos *matrix.Vector3 = &matrix.Vector3{X: 793, Y: 69, Z: 683}

var moveTarPosArr []*matrix.Vector3 = []*matrix.Vector3{
	&matrix.Vector3{X: 779, Y: 72, Z: 676},
	&matrix.Vector3{X: 766, Y: 75, Z: 665},
	&matrix.Vector3{X: 765, Y: 75.5, Z: 650},
	&matrix.Vector3{X: 789, Y: 73.5, Z: 655}, // ---
	&matrix.Vector3{X: 762, Y: 72, Z: 682},
	&matrix.Vector3{X: 745, Y: 76, Z: 667},
}

type MoveModel struct {
	net    *client_net.ClientNet
	userId int64

	speedMs float32
	mapId   int32
	pos     *matrix.Vector3
	dir     *matrix.Vector3

	preMoveMs int64
	targetPos *matrix.Vector3

	preSendMsgMs int64
}

func (m *MoveModel) SetUserId(userId int64)           { m.userId = userId }
func (m *MoveModel) SetNet(net *client_net.ClientNet) { m.net = net }
func (m *MoveModel) SetSpeed(spd float32)             { m.speedMs = spd / 1000.0 }
func (m *MoveModel) SetMapId(mapId int32)             { m.mapId = mapId }
func (m *MoveModel) SetPos(pos *matrix.Vector3)       { m.pos = pos }
func (m *MoveModel) SetDir(dir *matrix.Vector3)       { m.dir = dir }
func (m *MoveModel) IsMoving() bool                   { return m.targetPos != nil }

func (m *MoveModel) Tick(curMs int64) {
	if !m.IsMoving() {
		m.start(curMs)
		return
	}
	m.SendMoveMsg(curMs)
	m.move(curMs)
	m.stop(curMs)
}

func (m *MoveModel) getTargetPos() *matrix.Vector3 {
	canMoveIdx := make([]int, 0, len(moveTarPosArr))
	for idx, pos := range moveTarPosArr {
		if m.pos.Equal(*pos) {
			continue
		}
		canMoveIdx = append(canMoveIdx, idx)
	}

	idx := random.Random(0, len(canMoveIdx)-1)
	return moveTarPosArr[idx]

	// return &matrix.Vector3{X: 790, Y: 70, Z: 683}
}

func (m *MoveModel) start(curMs int64) {
	tarPos := m.getTargetPos()
	if tarPos.Equal(*m.pos) {
		return
	}

	m.targetPos = tarPos
	m.preMoveMs = curMs
	moveDir := matrix.Normalize3(matrix.Sub3(*m.targetPos, *m.pos))
	m.dir = &moveDir
	m.SendMoveMsg(curMs)
}

func (m *MoveModel) move(curMs int64) {
	if !m.IsMoving() {
		return
	}
	if m.pos.Equal(*m.targetPos) {
		return
	}

	moveDist := m.speedMs * float32(curMs-m.preMoveMs)
	toTargetDist := matrix.Distance(*m.pos, *m.targetPos)

	// 完成移动
	if toTargetDist <= float64(moveDist) {
		m.pos = m.targetPos
		return
	}

	// 移动
	moved := matrix.Multiply(*m.dir, float64(moveDist))
	toPos := matrix.Add3(*m.pos, moved)
	m.pos = &toPos
	m.preMoveMs = curMs
}

func (m *MoveModel) stop(curMs int64) {
	if !m.IsMoving() {
		return
	}
	if !m.pos.Equal(*m.targetPos) {
		return
	}

	m.preMoveMs = 0
	m.targetPos = nil
	m.SendStopMoveMsg(curMs)
	m.preSendMsgMs = 0
}

func (m *MoveModel) getCurMoveStep(curMs int64) *proto.EntityMoveStep {
	return &proto.EntityMoveStep{
		Stamp: curMs,
		Location: &proto.EntityLocation{
			MapId: m.mapId,
			Loc: &proto.Vector3{
				X: float32(m.pos.X), Y: float32(m.pos.Y), Z: float32(m.pos.Z),
			},
		},
	}
}

func (m *MoveModel) SendMoveMsg(curMs int64) {
	if m.targetPos == nil || curMs-m.preSendMsgMs < 300 {
		return
	}

	afterToMs := int64(500)
	moveDist := m.speedMs * float32(afterToMs)
	moved := matrix.Multiply(*m.dir, float64(moveDist))
	toPos := matrix.Add3(*m.pos, moved)

	destLocation := &proto.EntityMoveStep{
		Stamp: curMs + afterToMs,
		Location: &proto.EntityLocation{
			MapId: m.mapId,
			Loc: &proto.Vector3{
				X: float32(toPos.X), Y: float32(toPos.Y), Z: float32(toPos.Z),
			},
		},
	}

	movement := &proto.EntityMovement{
		TypeId: &proto.EntityId{
			Type: proto.EntityType_EntityTypePlayer,
			Id:   m.userId,
		},
		CurLocation:  m.getCurMoveStep(curMs),
		DestLocation: destLocation,
		Dir: &proto.Vector3{
			X: float32(m.dir.X), Y: float32(m.dir.Y), Z: float32(m.dir.Z),
		},
	}

	reqMsg := &proto.Envelope{
		Type: proto.EnvelopeType_UpdateSelfLocation,
		Payload: &proto.Envelope_UpdateSelfLocationRequest{
			UpdateSelfLocationRequest: &proto.UpdateSelfLocationRequest{
				Movement: movement,
			},
		},
	}
	m.net.Send(reqMsg)
	m.preSendMsgMs = curMs

	// serviceLog.Debug("[%d] sendMoveMsg  data: %+v", m.userId, movement)
}

func (m *MoveModel) SendStopMoveMsg(curMs int64) {
	movement := &proto.EntityMovement{
		TypeId: &proto.EntityId{
			Type: proto.EntityType_EntityTypePlayer,
			Id:   m.userId,
		},
		CurLocation:  m.getCurMoveStep(curMs),
		DestLocation: nil,
		Dir: &proto.Vector3{
			X: float32(m.dir.X), Y: float32(m.dir.Y), Z: float32(m.dir.Z),
		},
	}

	reqMsg := &proto.Envelope{
		Type: proto.EnvelopeType_UpdateSelfLocation,
		Payload: &proto.Envelope_UpdateSelfLocationRequest{
			UpdateSelfLocationRequest: &proto.UpdateSelfLocationRequest{
				Movement: movement,
			},
		},
	}
	m.net.Send(reqMsg)
	m.preSendMsgMs = curMs
	// serviceLog.Debug("[%d] STOP Move Msg  data: %+v", m.userId, movement)
}

func (m *MoveModel) OnUpdateSelfLocationRes(msg *proto.Envelope) {

}
