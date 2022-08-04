package weasel

import gonanoid "github.com/matoous/go-nanoid"

type Session interface {
	ReceiveWriter
	Serv
	SetKeeper(k *LRUKeeper)
	SessionId() string
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
}

type networkClient struct {
	keeper     *LRUKeeper
	MsgWriter  chan []byte
	MsgReader  chan []byte
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
		sessionId:  gonanoid.MustID(32),
		serialNo:   serialNo,
		serialName: serialName,
	}
}

func (p *networkClient) SerialNo() string {
	return p.serialNo
}

func (p *networkClient) SessionId() string {
	return p.sessionId
}

func (p *networkClient) SetKeeper(k *LRUKeeper) {
	p.keeper = k
}
