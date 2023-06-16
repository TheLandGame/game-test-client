package client

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/meland-client/src/client/client_ai"
	"github.com/Meland-Inc/meland-client/src/client/client_net"
	"github.com/Meland-Inc/meland-client/src/client/client_ping"
	"github.com/Meland-Inc/meland-client/src/client/data_model"
	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
)

type TestModel string

const (
	TEST_MODE_NORMAL    = "normal"
	TEST_MODE_CONNECT   = "connect"
	TEST_MODE_PING      = "ping"
	TEST_MODE_MAIN_DATA = "main-data"
)

type GameClient struct {
	model     TestModel
	net       client_net.ClientNet
	pingModel client_ping.ClientPing
	client_ai.ClientAiModel
	data_model.UserDataModel

	isStop  bool
	userIdx int64
	token   string

	msgEvent   map[proto.EnvelopeType]func(*net_packet.NetPacket)
	serMsgChan chan *net_packet.NetPacket
}

func NewGameClient(testModel string, agentUrl, token string, userIdx int64) *GameClient {
	c := &GameClient{
		userIdx:    userIdx,
		token:      token,
		msgEvent:   make(map[proto.EnvelopeType]func(*net_packet.NetPacket)),
		serMsgChan: make(chan *net_packet.NetPacket, 512),
	}
	if token == "" {
		c.token = fmt.Sprint(userIdx)
	}

	switch testModel {
	case string(TEST_MODE_CONNECT):
		c.model = TEST_MODE_CONNECT
	case string(TEST_MODE_PING):
		c.model = TEST_MODE_PING
	case string(TEST_MODE_MAIN_DATA):
		c.model = TEST_MODE_MAIN_DATA
	default:
		c.model = TEST_MODE_NORMAL
	}

	c.InitMsgHandler()
	c.net.Init(agentUrl, c.token, c.userIdx, c.MsgCallBack)
	c.pingModel.Init(&c.net)
	c.ClientAiModel.SetState(client_ai.USER_STATE_READY)
	c.UserDataModel.Init(&c.net)
	return c
}

func (c *GameClient) Run() {
	c.net.Start()
	c.start()

	for {
		c.tick(time_helper.NowUTCMill())
		if len(c.serMsgChan) <= 0 {
			time.Sleep(time.Millisecond * 1)
		}
	}
}

func (c *GameClient) tick(curMs int64) {
	if c.isStop {
		return
	}
	// defer func ()  {
	// 	if err := recover();err!=nil{
	// 		serviceLog.StackError()
	// 	}
	// }()

	c.readSerMsg()
	c.pingTick(curMs)
	c.logicTick(curMs)
}

func (c *GameClient) MsgCallBack(packet *net_packet.NetPacket) {
	// if len(c.serMsgChan) > 5 && c.userIdx%10000 == 0 {
	// 	serviceLog.Warning("userIdx:%v serMsgChan len %d  > 5", c.userIdx, len(c.serMsgChan))
	// }

	c.serMsgChan <- packet
}

func (c *GameClient) readSerMsg() {
	select {
	case packet := <-c.serMsgChan:
		c.onReceiveMsg(packet)
	default:
	}
}

func (c *GameClient) start() {
	if c.model == TEST_MODE_CONNECT || c.model == TEST_MODE_PING {
		return
	}
	c.QueryUser()
}

func (c *GameClient) stop() {
	c.isStop = true
	c.net.Stop()
}

func (c *GameClient) pingTick(curMs int64) {
	if c.model == TEST_MODE_CONNECT {
		return
	}
	c.pingModel.Tick(curMs)
}

func (c *GameClient) logicTick(curMs int64) {
	switch c.model {
	case TEST_MODE_NORMAL:
		c.ClientAiModel.Tick(curMs)
	case TEST_MODE_MAIN_DATA:
		c.UserDataModel.Tick(curMs)
	}
}
