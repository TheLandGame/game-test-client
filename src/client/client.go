package client

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/meland-client/src/client/client_ai"
	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/client/client_ping"
	"github.com/Meland-Inc/meland-client/src/common/matrix"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

type TestModel string

const (
	TEST_MODE_NORMAL  = "normal"
	TEST_MODE_CONNECT = "connect"
	TEST_MODE_PING    = "ping"
)

type UserData struct {
	baseData  *proto.PlayerBaseData
	sceneData *proto.Player
	mapId     int32
	Pos       *matrix.Vector3
	Dir       *matrix.Vector3
}

type GameClient struct {
	model     TestModel
	net       client_net.ClientNet
	pingModel client_ping.ClientPing
	client_ai.ClientAiModel

	userIdx int64
	token   string

	msgEvent   map[proto.EnvelopeType]func(*proto.Envelope)
	serMsgChan chan *proto.Envelope

	playerData UserData
}

func NewGameClient(testModel string, agentUrl, token string, userIdx int64) *GameClient {
	c := &GameClient{
		userIdx:    userIdx,
		token:      token,
		msgEvent:   make(map[proto.EnvelopeType]func(*proto.Envelope)),
		serMsgChan: make(chan *proto.Envelope, 512),
	}
	if token == "" {
		c.token = fmt.Sprint(userIdx)
	}

	switch testModel {
	case string(TEST_MODE_CONNECT):
		c.model = TEST_MODE_CONNECT
	case string(TEST_MODE_PING):
		c.model = TEST_MODE_PING
	default:
		c.model = TEST_MODE_NORMAL
	}

	c.ClientAiModel.SetState(client_ai.USER_STATE_READY)

	c.InitMsgHandler()
	c.net.Init(agentUrl, c.token, c.userIdx, c.MsgCallBack)
	c.pingModel.Init(&c.net)
	return c
}

func (c *GameClient) Run() {
	c.net.Start()
	c.start()

	for {
		c.tick(time_helper.NowUTCMill())
		time.Sleep(time.Millisecond * 10)
	}
}

func (c *GameClient) registerMsgHandler(
	msgType proto.EnvelopeType, handler func(*proto.Envelope),
) {
	c.msgEvent[msgType] = handler
}

func (c *GameClient) MsgCallBack(msg *proto.Envelope) {
	c.serMsgChan <- msg
}

func (c *GameClient) readMs() {
	select {
	case msg := <-c.serMsgChan:
		c.onReceiveMsg(msg)
	default:
	}
}

func (c *GameClient) start() {
	if c.model != TEST_MODE_NORMAL {
		return
	}
	c.QueryUser()
}

func (c *GameClient) tick(curMs int64) {
	c.readMs()
	c.pingTick(curMs)
	c.logicTick(curMs)
}

func (c *GameClient) pingTick(curMs int64) {
	if c.model == TEST_MODE_CONNECT {
		return
	}
	c.pingModel.Tick(curMs)
}

func (c *GameClient) logicTick(curMs int64) {
	if c.model != TEST_MODE_NORMAL {
		return
	}
	c.ClientAiModel.Tick(curMs)
}
