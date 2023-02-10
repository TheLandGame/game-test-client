package client

import (
	"game-message-core/proto"
	"log"

	"github.com/Meland-Inc/meland-client/src/client/msg_record"
	"github.com/Meland-Inc/meland-client/src/common/time_helper"
	"github.com/Meland-Inc/meland-client/src/net_util"
)

func (c *GameClient) PrintMsgUsedMs(respMsg *proto.Envelope) {
	iRecord, exist := c.callback.LoadAndDelete(respMsg.SeqId)
	if !exist {
		return
	}
	curMs := time_helper.NowUTCMicro()
	record := iRecord.(*msg_record.MsgRecord)
	log.Printf(
		"[I]: Msg[%v] SeaId[%d], reqUs[%d], respUs[%d],  通信耗时[%5f]MS\n",
		respMsg.Type, respMsg.SeqId, record.CreateAt, curMs, float32(curMs-record.CreateAt)/1000,
	)
}

func (c *GameClient) onData(data []byte) {
	msg, err := net_util.UnMarshalProtoMessage(data)
	if err != nil {
		panic(err)
	}

	c.PrintMsgUsedMs(msg)

	switch msg.Type {
	case proto.EnvelopeType_Ping:
		c.PingHandler(msg)

	}
}

func (c *GameClient) PingHandler(msg *proto.Envelope) {

}
