package msg_record

import "game-message-core/proto"

type MsgRecord struct {
	CreateAt int64
	Type     proto.EnvelopeType
}
