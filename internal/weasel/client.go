package weasel

import gonanoid "github.com/matoous/go-nanoid"

type Session interface {
	ReceiveWriter
	Serv
	Close()
}

type ReceiveWriter interface {
	Write(b []byte)
	Receive() <-chan []byte
}

type Serv interface {
	SerialNo() string
	WriterServ()
	ReaderServ()
	Dead() <-chan bool
	SessionId() string
}

type networkClient struct {
	MsgWriter chan []byte
	MsgReader chan []byte

	serialNo   string
	serialName string
	sessionId  string
}

//	an alias to point networkClient
type NetworkClient = *networkClient

func NewNetworkClient(serialNo, serialName string) NetworkClient {
	return &networkClient{
		MsgWriter:  make(chan []byte),
		MsgReader:  make(chan []byte),
		serialNo:   serialNo,
		serialName: serialName,
		sessionId:  gonanoid.MustID(32),
	}
}

func (p *networkClient) SerialNo() string {
	return p.serialNo
}

func (p *networkClient) SessionId() string {
	return p.sessionId
}
