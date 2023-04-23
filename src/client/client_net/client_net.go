package client_net

import (
	"fmt"
	"game-message-core/proto"
	"sync"

	"github.com/Meland-Inc/meland-client/src/client/msg_record"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
	"github.com/Meland-Inc/meland-client/src/net_util"
	"github.com/gorilla/websocket"
)

type ClientNet struct {
	agentUrl    string
	token       string
	userId      int64
	MsgCallBack func(*proto.Envelope)

	wcConnect *websocket.Conn
	msgSeqId  int32

	reqMsgRecord sync.Map //{msgSeqId = *MsgRecord}
}

func (c *ClientNet) Init(agentUrl, token string, userId int64, msgCallBack func(msg *proto.Envelope)) {
	if agentUrl == "" || token == "" {
		panic("invalid agentUrl  ||  token")
	}
	if msgCallBack == nil {
		panic("invalid msgCallBack")
	}

	c.agentUrl = agentUrl
	c.token = token
	c.MsgCallBack = msgCallBack
	c.userId = userId
}

func (c *ClientNet) NextSeqId() int32 {
	c.msgSeqId++
	return c.msgSeqId
}

func (c *ClientNet) Start() {
	if err := c.connect(); err != nil {
		panic(err)
	}
	c.run()
}

func (c *ClientNet) connect() error {
	var dialer *websocket.Dialer
	url := fmt.Sprintf("ws://%s?token=%s", c.agentUrl, c.token)
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.wcConnect = conn
	serviceLog.Info("userIdx[%d] websocket connect success", c.userId)
	return nil
}

func (c *ClientNet) run() {
	go func() {
		for {
			_, data, err := c.wcConnect.ReadMessage()
			if err != nil {
				panic(err)
			}

			c.onData(data)
		}
	}()
}

func (c *ClientNet) onData(data []byte) {
	msg, err := net_util.UnMarshalProtoMessage(data)
	if err != nil {
		panic(err)
	}

	c.PrintMsgUsedMs(msg)
	c.MsgCallBack(msg)
}

func (c *ClientNet) Send(msg *proto.Envelope) {
	msg.SeqId = c.NextSeqId()
	payload, err := net_util.MarshalProtoMessage(msg)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	c.onSendMsg(msg)
	err = c.wcConnect.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		panic(err)
	}
}

func (c *ClientNet) onSendMsg(msg *proto.Envelope) {
	if msg.Type == proto.EnvelopeType_Ping {
		return
	}

	// serviceLog.Debug("cli[%d] send msg [%v], seqId[%d]", c.userId, msg.Type, msg.SeqId)

	c.reqMsgRecord.Store(msg.SeqId, &msg_record.MsgRecord{
		CreateAt: time_helper.NowUTCMill(),
		ReqMsg:   msg,
	})
}

func (c *ClientNet) PrintMsgUsedMs(respMsg *proto.Envelope) {
	if c.userId != 20000 {
		return
	}

	curMs := time_helper.NowUTCMill()
	// serviceLog.Debug(
	// 	"[I]: %d Msg[%v] SeaId[%d],  \n", curMs, respMsg.Type, respMsg.SeqId,
	// )

	if respMsg.Type == proto.EnvelopeType_BroadCastMsgAggregation {
		res := respMsg.GetBroadCastMsgAggregationResponse()
		for _, aMsg := range res.MessageList {
			c.PrintMsgUsedMs(aMsg)
		}
		return
	}

	iRecord, exist := c.reqMsgRecord.LoadAndDelete(respMsg.SeqId)
	if !exist {
		return
	}
	record := iRecord.(*msg_record.MsgRecord)
	serviceLog.Info(
		"user[%d], Msg[%v] SeqId[%d], reqUs[%d], respUs[%d],  通信耗时[%d]MS",
		c.userId, respMsg.Type, respMsg.SeqId, record.CreateAt, curMs, curMs-record.CreateAt,
	)
}
