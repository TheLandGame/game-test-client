package client_ai

import (
	"github.com/Meland-Inc/meland-client/src/client/client_ai/move_model"
	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
	"github.com/Meland-Inc/meland-client/src/common/random"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

type ClientAiModel struct {
	move_model.MoveModel
	userId int64
	net    *client_net.ClientNet
	state  UserState

	startMs     int64
	startIdleMs int64
}

func (ai *ClientAiModel) Init(
	net *client_net.ClientNet, userId int64, mapId int32, pos, dir *matrix.Vector3, userSpd float32,
) {
	ai.net = net
	ai.MoveModel.SetNet(net)
	ai.userId = userId
	ai.MoveModel.SetUserId(userId)
	ai.MoveModel.SetSpeed(userSpd)
	ai.MoveModel.SetMapId(mapId)
	ai.MoveModel.SetPos(pos)
	ai.MoveModel.SetDir(dir)
	ai.startIdleMs = time_helper.NowUTCMill()
	ai.state = USER_STATE_IDLE
	// serviceLog.Info("user[%d] ClientAi init SUCCESS ~_~ ", ai.userId)
}

func (ai *ClientAiModel) SetState(state UserState) {
	ai.state = state
}

func (ai *ClientAiModel) Tick(curMs int64) {
	switch ai.state {
	case USER_STATE_READY:
		// if curMs-ai.startMs > 3000 {
		// 	ai.state = USER_STATE_IDLE
		// 	ai.startIdleMs = curMs
		// 	serviceLog.Debug("user[%d] USER_STATE_READY  ->  USER_STATE_IDLE", ai.userId)
		// }

	case USER_STATE_IDLE:
		if curMs-ai.startIdleMs < 3000 {
			return
		}
		rdn := random.Random32(1, 100)
		if rdn < 10 {
			ai.state = USER_STATE_MOVE
			// serviceLog.Debug("user[%d] USER_STATE_IDLE  ->  USER_STATE_MOVE", ai.userId)
		}

	case USER_STATE_MOVE:
		ai.MoveModel.Tick(curMs)
		if !ai.MoveModel.IsMoving() {
			ai.state = USER_STATE_IDLE
			ai.startIdleMs = curMs
			// serviceLog.Debug("user[%d] USER_STATE_MOVE  ->  USER_STATE_IDLE", ai.userId)
		}

	case USER_STATE_ATTACK:

	default:
		serviceLog.Debug("user[%d]  --------------------------       def =======", ai.userId)
	}
}
