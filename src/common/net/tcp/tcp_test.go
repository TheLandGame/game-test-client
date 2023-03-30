package tcp

import (
	"fmt"
	"game-message-core/proto"
	"game-message-core/protoTool"
	"testing"
	"time"

	"github.com/Meland-Inc/meland-client/src/common/net/session"
)

func ServerOnConnectCallback(s *session.Session) {
	fmt.Printf("session[%s][%s] ServerOnConnectCallback  ------\n", s.SessionId(), s.RemoteAddr())
	s.SetCallBack(
		ServerOnReceivedCallback,
		ServerOnCloseCallback,
	)
}
func ServerOnReceivedCallback(s *session.Session, data []byte) {
	fmt.Printf(
		"session[%s][%s]  ServerOnReceivedCallback data length[%d], dataStr:[%v]\n",
		s.SessionId(), s.RemoteAddr(), len(data), string(data),
	)
	s.Write([]byte("server response message is ServerOnReceivedCallback ########"))
}
func ServerOnCloseCallback(s *session.Session) {
	fmt.Printf("session[%s][%s]  ServerOnCloseCallback ======= \n", s.SessionId(), s.RemoteAddr())
}

func Test_TcpServer(t *testing.T) {
	tcpServer, err := NewTcpServer(
		":7659",
		100,
		180,
		ServerOnConnectCallback,
	)
	fmt.Println(err)
	fmt.Println(tcpServer)
	fmt.Printf("tcp server started  \n")
	time.Sleep(1 * time.Hour)
}

func ClientOnReceivedCallback(s *session.Session, data []byte) {
	fmt.Printf(
		"session[%s][%s]  ClientOnReceivedCallback data length[%d], dataStr:[%v]\n",
		s.SessionId(), s.RemoteAddr(), len(data), string(data),
	)
}
func ClientOnCloseCallback(s *session.Session) {
	fmt.Printf("session[%s][%s]  ClientOnCloseCallback ======= \n", s.SessionId(), s.RemoteAddr())
}
func Test_TcpClient(t *testing.T) {
	t.Log("TcpClient ------- begin --------")
	cli, err := NewClient(":5700", ClientOnReceivedCallback, ServerOnCloseCallback)
	if err != nil {
		panic(err)
	}

	t.Log("TcpClient ------- send --------")
	msg := &proto.Envelope{
		SeqId: 44444,
		Type:  proto.EnvelopeType_Ping,
		Payload: &proto.Envelope_PingRequest{
			PingRequest: &proto.PingRequest{},
		},
	}
	bs, err := protoTool.MarshalProto(msg)
	t.Log(err)
	t.Log(bs)

	cli.Write(bs)

	t.Log("TcpClient ------- send over --------")
	time.Sleep(10 * time.Second)
	cli.Close()
}
