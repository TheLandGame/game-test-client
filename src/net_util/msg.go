package net_util

import (
	"game-message-core/proto"
)

func MakeEnterMapMsg() *proto.EnterMapReq {
	return &proto.EnterMapReq{
		ReqTitle: &proto.ReqTitle{},
	}
}

func MakeItemGetMsg() *proto.ItemGetReq {
	return &proto.ItemGetReq{
		ReqTitle: &proto.ReqTitle{},
	}
}
