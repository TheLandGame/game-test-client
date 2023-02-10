package net_util

import (
	"game-message-core/proto"
)

func MakePingMsg() *proto.Envelope {
	return &proto.Envelope{
		Type: proto.EnvelopeType_Ping,
		Payload: &proto.Envelope_PingRequest{
			PingRequest: &proto.PingRequest{},
		},
	}
}

func MakeQueryPlayerMsg() *proto.Envelope {
	return &proto.Envelope{
		Type: proto.EnvelopeType_QueryPlayer,
		Payload: &proto.Envelope_QueryPlayerRequest{
			QueryPlayerRequest: &proto.QueryPlayerRequest{
				// Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inc4NWN0ZWQiLCJzdWIiOiJ7XCJpZFwiOlwiNjgwXCIsXCJzY2hvb2xJZFwiOlwiMFwiLFwidXNlcm5hbWVcIjpcInc4NWN0ZWRcIixcIm5pY2tuYW1lXCI6XCJ3MTU5ODY3N1wiLFwicmVhbG5hbWVcIjpcIncxNTk4Njc3XCIsXCJ1c2VydHlwZVwiOlwiU1RVREVOVFwiLFwic2V4XCI6XCJNQUxFXCIsXCJlbWFpbFwiOlwiXCIsXCJtb2JpbGVcIjpcIlwiLFwiYXZhdGFyXCI6XCJodHRwczovL2JnLXByby53a2NvZGluZy5jb20vZGVmYXVsdF9hdmF0YXItMjAyMTEwMjAucG5nXCIsXCJyZWdpc3RlZEF0XCI6XCIyMDIyLTA2LTMwVDA3OjExOjIwLjAwMFpcIixcImxhc3RMb2dpbmVkQXRcIjpcIjIwMjItMDYtMzBUMDc6MTE6MjAuMDAwWlwiLFwibGFzdFVwZGF0ZWRBdFwiOlwiMjAyMi0wNi0zMFQwNzoxMToyMC4wMDBaXCIsXCJzdGF0dXNcIjoxfSIsImlhdCI6MTY2MzIyODU4NywiZXhwIjoxNjYzODMzMzg3fQ.RMn62oYTMiReocxlnZ5Q149DvFyi2vVmfLRAtolcHqc",
			},
		},
	}
}

func MakeCreatePlayerMsg() *proto.Envelope {
	return &proto.Envelope{
		Type: proto.EnvelopeType_CreatePlayer,
		Payload: &proto.Envelope_CreatePlayerRequest{
			CreatePlayerRequest: &proto.CreatePlayerRequest{
				// Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inc4NWN0ZWQiLCJzdWIiOiJ7XCJpZFwiOlwiNjgwXCIsXCJzY2hvb2xJZFwiOlwiMFwiLFwidXNlcm5hbWVcIjpcInc4NWN0ZWRcIixcIm5pY2tuYW1lXCI6XCJ3MTU5ODY3N1wiLFwicmVhbG5hbWVcIjpcIncxNTk4Njc3XCIsXCJ1c2VydHlwZVwiOlwiU1RVREVOVFwiLFwic2V4XCI6XCJNQUxFXCIsXCJlbWFpbFwiOlwiXCIsXCJtb2JpbGVcIjpcIlwiLFwiYXZhdGFyXCI6XCJodHRwczovL2JnLXByby53a2NvZGluZy5jb20vZGVmYXVsdF9hdmF0YXItMjAyMTEwMjAucG5nXCIsXCJyZWdpc3RlZEF0XCI6XCIyMDIyLTA2LTMwVDA3OjExOjIwLjAwMFpcIixcImxhc3RMb2dpbmVkQXRcIjpcIjIwMjItMDYtMzBUMDc6MTE6MjAuMDAwWlwiLFwibGFzdFVwZGF0ZWRBdFwiOlwiMjAyMi0wNi0zMFQwNzoxMToyMC4wMDBaXCIsXCJzdGF0dXNcIjoxfSIsImlhdCI6MTY2MzIyODU4NywiZXhwIjoxNjYzODMzMzg3fQ.RMn62oYTMiReocxlnZ5Q149DvFyi2vVmfLRAtolcHqc",
				// NickName: "黄玥玥",
				// RoleId:   1002,
				// Gender:   "man",
				// Icon:     "icon---image",
				// Feature: &proto.PlayerFeature{
				// 	Face:    100003,
				// 	Hair:    101002,
				// 	Glove:   105001,
				// 	Clothes: 102001,
				// 	Pants:   103001,
				// },
			},
		},
	}
	//{"hair":101002,"clothes":102001,"glove":105001,"pants":103001,"face":100003,"shoes":104002}
}

func makeSingInMsg() *proto.Envelope {
	return &proto.Envelope{
		Type: proto.EnvelopeType_SigninPlayer,
		Payload: &proto.Envelope_SigninPlayerRequest{
			SigninPlayerRequest: &proto.SigninPlayerRequest{
				// Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inc4NWN0ZWQiLCJzdWIiOiJ7XCJpZFwiOlwiNjgwXCIsXCJzY2hvb2xJZFwiOlwiMFwiLFwidXNlcm5hbWVcIjpcInc4NWN0ZWRcIixcIm5pY2tuYW1lXCI6XCJ3MTU5ODY3N1wiLFwicmVhbG5hbWVcIjpcIncxNTk4Njc3XCIsXCJ1c2VydHlwZVwiOlwiU1RVREVOVFwiLFwic2V4XCI6XCJNQUxFXCIsXCJlbWFpbFwiOlwiXCIsXCJtb2JpbGVcIjpcIlwiLFwiYXZhdGFyXCI6XCJodHRwczovL2JnLXByby53a2NvZGluZy5jb20vZGVmYXVsdF9hdmF0YXItMjAyMTEwMjAucG5nXCIsXCJyZWdpc3RlZEF0XCI6XCIyMDIyLTA2LTMwVDA3OjExOjIwLjAwMFpcIixcImxhc3RMb2dpbmVkQXRcIjpcIjIwMjItMDYtMzBUMDc6MTE6MjAuMDAwWlwiLFwibGFzdFVwZGF0ZWRBdFwiOlwiMjAyMi0wNi0zMFQwNzoxMToyMC4wMDBaXCIsXCJzdGF0dXNcIjoxfSIsImlhdCI6MTY2MzIyODU4NywiZXhwIjoxNjYzODMzMzg3fQ.RMn62oYTMiReocxlnZ5Q149DvFyi2vVmfLRAtolcHqc",
				// ClientTime: time_helper.NowUTCMill(),
				// Reconnect:  false,
				// SceneServiceAppId: "scene_service_scene_801",
			},
		},
	}
}

func makeItemGetMsg() *proto.Envelope {
	return &proto.Envelope{
		Type: proto.EnvelopeType_ItemGet,
		Payload: &proto.Envelope_ItemGetRequest{
			ItemGetRequest: &proto.ItemGetRequest{},
		},
	}
}

func makeEnterMapMsg() *proto.Envelope {
	return &proto.Envelope{
		Type: proto.EnvelopeType_EnterMap,
		Payload: &proto.Envelope_EnterMapRequest{
			EnterMapRequest: &proto.EnterMapRequest{},
		},
	}
}
