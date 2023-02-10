package net_util

import (
	"game-message-core/proto"
	"game-message-core/protoTool"
)

func UnMarshalProtoMessage(data []byte) (*proto.Envelope, error) {
	msg := &proto.Envelope{}
	err := protoTool.UnmarshalProto(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func MarshalProtoMessage(msg *proto.Envelope) ([]byte, error) {
	bs, err := protoTool.MarshalProto(msg)
	if err != nil {
		return nil, err
	}
	return bs, err
}
