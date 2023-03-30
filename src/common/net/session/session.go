package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Meland-Inc/meland-client/src/common/net/msgParser"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/google/uuid"
)

const (
	WRITE_CHAN_SIZE = 2000
)

type Session struct {
	id         string
	owner      int32
	conn       net.Conn
	reader     *bufio.Reader
	parser     *msgParser.MsgParser
	activeTime int64
	closed     bool

	sendChan           chan []byte
	stopChan           chan chan struct{}
	onReceivedCallback func(*Session, []byte)
	onCloseCallback    func(*Session)
}

func NewSession(con net.Conn) *Session {
	session := &Session{
		id:         uuid.New().String(),
		conn:       con,
		reader:     bufio.NewReader(con),
		activeTime: time.Now().UTC().Unix(),
		parser:     msgParser.NewMsgParser(),
		sendChan:   make(chan []byte, WRITE_CHAN_SIZE),
		stopChan:   make(chan chan struct{}),
	}
	return session
}

func (s *Session) SessionId() string { return s.id }

func (s *Session) SetOwner(ownerId int32) { s.owner = ownerId }

func (s *Session) GetOwner() int32 { return s.owner }

func (s *Session) IsClosed() bool { return s.closed }

func (s *Session) LocalAddr() net.Addr { return s.conn.LocalAddr() }

func (s *Session) RemoteAddr() net.Addr { return s.conn.RemoteAddr() }

func (s *Session) SetCallBack(
	onDataCallback func(*Session, []byte),
	onCloseCallback func(*Session),
) error {
	if onCloseCallback == nil || onDataCallback == nil {
		return fmt.Errorf("invalid call back function")
	}
	s.onReceivedCallback = onDataCallback
	s.onCloseCallback = onCloseCallback
	return nil
}

func (s *Session) GetActiveTime() int64 { return s.activeTime }

func (s *Session) SetActiveTime(activeTime int64) {
	if activeTime <= s.activeTime {
		return
	}
	s.activeTime = activeTime
}

func (s *Session) String() string {
	if s.conn != nil {
		return "[" + s.conn.RemoteAddr().String() + "]"
	} else {
		out, err := json.Marshal(s)
		if err != nil {
			return err.Error()
		}

		return string(out)
	}
}

func (s *Session) Run() {
	s.closed = false

	go func() {
		s.loop()
	}()

	err := s.received()
	if err != nil {
		serviceLog.Warning("session [%v] received err: %v", s.RemoteAddr(), err)
	}
}

func (s *Session) loop() {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.Error("session onSend err: ", err)
			go s.loop()
		}
	}()

	for {
		select {
		case msg := <-s.sendChan:
			if err := s.send(msg); err != nil {
				serviceLog.Warning("Stop session(%d)  by send message err :  %v", s.RemoteAddr(), err)
				go s.Stop()
			}

		case stopFinished := <-s.stopChan:
			serviceLog.Info("Stop session(%v)  by stop event", s.RemoteAddr())
			s.close()
			close(s.sendChan)
			close(s.stopChan)
			stopFinished <- struct{}{}
			return
		}
	}
}

func (s *Session) received() error {
	for {
		if s.IsClosed() {
			break
		}

		head := make([]byte, msgParser.HEAD_SIZE)
		data, err := s.parser.Decode(head, s.reader)
		if err != nil {
			return err
		}
		s.onReceivedCallback(s, data)
		s.SetActiveTime(time.Now().UTC().Unix())
	}
	return nil
}

func (s *Session) send(data []byte) error {
	if s.IsClosed() {
		return fmt.Errorf("Session closed!! by sendChan")
	}

	bytes, err := s.parser.Encode(data)
	if err != nil {
		return err
	}

	_, err = s.conn.Write(bytes)
	return err
}

func (s *Session) Write(data []byte) error {
	if s.IsClosed() {
		return fmt.Errorf("Session closed!! by sendChan")
	}
	if len(data) > msgParser.MSG_LIMIT {
		return fmt.Errorf("write msg length exceeds msgParser.MSG_LIMIT[%v]", msgParser.MSG_LIMIT)
	}
	s.sendChan <- data
	return nil
}

func (s *Session) close() {
	if !s.IsClosed() {
		s.conn.(*net.TCPConn).SetLinger(0)
		s.conn.Close()
		s.closed = true
		if s.onCloseCallback != nil {
			s.onCloseCallback(s)
		}
	}
	s.closed = true
}

func (s *Session) Stop() {
	if !s.IsClosed() {
		stopDone := make(chan struct{}, 1)
		s.stopChan <- stopDone
		<-stopDone
	}
}
