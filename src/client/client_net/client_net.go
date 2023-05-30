package client_net

import (
	"fmt"
	"game-message-core/proto"
	"game-message-core/protoTool"
	"sync"
	"time"

	"github.com/Meland-Inc/meland-client/src/client/msg_record"
	"github.com/Meland-Inc/meland-client/src/common/net/net_packet"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
	"github.com/gorilla/websocket"
	googleProto "google.golang.org/protobuf/proto"
)

type ClientNet struct {
	agentUrl    string
	token       string
	userId      int64
	MsgCallBack func(*net_packet.NetPacket)

	wcConnect *websocket.Conn
	msgSeqId  int32

	reqMsgRecord sync.Map //{msgSeqId = *MsgRecord}
}

func (c *ClientNet) Init(
	agentUrl, token string, userId int64, msgCallBack func(msg *net_packet.NetPacket),
) {
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

func (c *ClientNet) Stop() {
	if c.wcConnect == nil {
		return
	}
	c.wcConnect.Close()
	c.wcConnect = nil
}
func (c *ClientNet) Start() {
	for {
		err := c.connect()
		if err == nil {
			c.run()
			return
		}
		serviceLog.Error("userIdx[%d] websocket connect failed  err: %v", c.userId, err.Error())
		time.Sleep(time.Millisecond * 200)
	}
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
				serviceLog.Error("userIdx[%d] websocket ReadMessage failed  err: %v", c.userId, err.Error())
				// panic(err)
				c.Stop()
				return
			}

			c.onData(data)
		}
	}()
}

func (c *ClientNet) onData(data []byte) {
	// serviceLog.Debug("receive data Len[%d] ", len(data))
	netPackets, err := net_packet.ReadPacket(data)
	if err != nil {
		serviceLog.Error("clientNet Ondata err: %s", err.Error())
		return
	}

	for _, packet := range netPackets {
		c.MsgCallBack(packet)
	}
}

func (c *ClientNet) Send(eType proto.EnvelopeType, msg googleProto.Message) {
	if c.wcConnect == nil {
		return
	}

	bs, err := protoTool.MarshalProto(msg)
	if err != nil {
		serviceLog.Error("Send Type:%v err:", eType, err.Error())
		return
	}

	packet := net_packet.GetPool().Get()
	packet.Id = uint32(eType)
	packet.Length = len(bs)
	packet.Body = bs

	bodyList := net_packet.WritePacket([]*net_packet.NetPacket{packet})
	for _, body := range bodyList {
		// serviceLog.Debug("Send [%d] body: %+v", eType, body)
		err = c.wcConnect.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			// panic(err)
			serviceLog.Error("userIdx[%d] eType[%v] WriteMessage failed  err: %v", c.userId, eType, err.Error())
		}
	}
}

// func (c *ClientNet) SendTestarr() {
// 	packets := []*net_packet.NetPacket{}
// 	for i := int32(1); i <= 5; i++ {
// 		req := &proto.PingReq{
// 			ReqTitle: &proto.ReqTitle{SeqId: i},
// 		}
// 		bs, err := protoTool.MarshalProto(req)
// 		if err != nil {
// 			serviceLog.Error(err.Error())
// 			return
// 		}
// 		packet := net_packet.GetPool().Get()
// 		packet.Id = uint32(proto.EnvelopeType_Ping)
// 		packet.Length = len(bs)
// 		packet.Body = bs
// 		packets = append(packets, packet)
// 	}
// 	bodyList := net_packet.WritePacket(packets)
// 	for _, body := range bodyList {
// 		// serviceLog.Debug("Send body: %+v", body)
// 		err := c.wcConnect.WriteMessage(websocket.BinaryMessage, body)
// 		if err != nil {
// 			// panic(err)
// 			serviceLog.Error("userIdx[%d] eType[%v] WriteMessage failed  err: %v", c.userId, err.Error())
// 		}
// 	}
// }

func (c *ClientNet) OnSendMsg(eType proto.EnvelopeType, seqId int32) {
	if c.userId%10000 != 0 {
		return
	}

	// if eType == proto.EnvelopeType_Ping {
	// 	return
	// }

	// serviceLog.Debug("cli[%d] send msg [%v], seqId[%d]", c.userId, eType, seqId)

	c.reqMsgRecord.Store(seqId, &msg_record.MsgRecord{
		CreateAt: time_helper.NowUTCMill(),
		Type:     eType,
	})
}

func (c *ClientNet) PrintMsgUsedMs(eType proto.EnvelopeType, seqId int32) {
	if c.userId%10000 != 0 {
		return
	}

	curMs := time_helper.NowUTCMill()

	iRecord, exist := c.reqMsgRecord.LoadAndDelete(seqId)
	if !exist {
		return
	}

	record := iRecord.(*msg_record.MsgRecord)
	serviceLog.Info(
		"user[%d], Msg[%v] SeqId[%d], reqUs[%d], respUs[%d],  通信耗时[%d]MS",
		c.userId, eType, seqId, record.CreateAt, curMs, curMs-record.CreateAt,
	)
}
