package client

import (
	"fmt"
	"game-message-core/proto"
	"log"
	"sync"
	"time"

	"github.com/Meland-Inc/meland-client/src/client/msg_record"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
	"github.com/Meland-Inc/meland-client/src/net_util"
	"github.com/gorilla/websocket"
)

type GameClient struct {
	agentUrl  string
	token     string
	wcConnect *websocket.Conn

	msgSeqId int32
	callback sync.Map //{msgSeqId = *MsgRecord}
}

func NewGameClient(agentUrl, token string) *GameClient {
	c := &GameClient{
		agentUrl: agentUrl,
		token:    token,
	}
	return c
}
func (c *GameClient) NextSeqId() int32 {
	c.msgSeqId++
	return c.msgSeqId
}

func (c *GameClient) Run() error {
	if err := c.connect(); err != nil {
		panic(err)
	}
	c.Ping()

	go func() {
		for {
			_, data, err := c.wcConnect.ReadMessage()
			if err != nil {
				panic(err)
			}
			c.onData(data)
		}
	}()

	return nil
}

func (c *GameClient) Ping() {
	msg := net_util.MakePingMsg()
	c.Send(msg)

	time.AfterFunc(time.Millisecond*500, func() {
		c.Ping()
	})
}

func (c *GameClient) connect() error {
	var dialer *websocket.Dialer
	url := fmt.Sprintf("ws://%s?token=%s", c.agentUrl, c.token)
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.wcConnect = conn
	return nil
}

func (c *GameClient) Send(msg *proto.Envelope) {
	msg.SeqId = c.NextSeqId()
	payload, err := net_util.MarshalProtoMessage(msg)
	if err != nil {
		log.Printf("[E]: %s \n", err.Error())
		return
	}

	log.Printf("[D]: client send msg [%v], seqId[%d] \n", msg.Type, msg.SeqId)
	c.onSendMsg(msg)
	err = c.wcConnect.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		panic(err)
	}
}

func (c *GameClient) onSendMsg(msg *proto.Envelope) {
	c.callback.Store(msg.SeqId, &msg_record.MsgRecord{
		CreateAt: time_helper.NowUTCMicro(),
		ReqMsg:   msg,
	})
}
