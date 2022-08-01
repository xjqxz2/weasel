package weasel

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
}

type networkClient struct {
	MsgWriter chan []byte
	MsgReader chan []byte

	serialNo   string
	serialName string
}

//	an alias to point networkClient
type NetworkClient = *networkClient

func NewNetworkClient(serialNo, serialName string) NetworkClient {
	return &networkClient{
		MsgWriter:  make(chan []byte),
		MsgReader:  make(chan []byte),
		serialNo:   serialNo,
		serialName: serialName,
	}
}

func (p *networkClient) SerialNo() string {
	return p.serialNo
}
