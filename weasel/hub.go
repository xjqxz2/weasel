package weasel

import (
	"sync"
)

type Hub struct {
	rmu     sync.RWMutex
	clients map[string]Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]Client),
	}
}

//	将客户端（连接）注册至 Hub 中
func (p *Hub) Register(serialNo string, client Client) error {
	if client := p.Find(serialNo); client != nil {
		if err := p.UnRegister(client); err != nil {
			return err
		}
	}

	p.rmu.Lock()
	defer p.rmu.Unlock()

	p.clients[serialNo] = client

	return nil
}

func (p *Hub) UnRegister(client Client) error {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	//	Close Old Connection
	client.Close()
	delete(p.clients, client.SerialNo())

	return nil
}

func (p *Hub) Find(serialNo string) Client {
	p.rmu.RLock()
	defer p.rmu.RUnlock()

	return p.clients[serialNo]
}

func (p *Hub) Search(serialsNo ...string) BroadcastTarget {
	var result BroadcastTarget

	//	当未指定广播设备序列时，则使用完全广播模式
	if len(serialsNo) <= 0 {
		for _, client := range p.clients {
			result = append(result, client)
		}
	}

	//	当指定广播设备序列时，则使用局部广播模式
	for _, serialNo := range serialsNo {
		if client, ok := p.clients[serialNo]; ok {
			result = append(result, client)
		}
	}

	return result
}

func (p *Hub) Start(client Client) {
	//	开始消息循环
	client.WriterServ()
	client.ReaderServ()
}
