package weasel

import (
	"log"
	"sync"
)

type Hub struct {
	rmu      sync.RWMutex
	sessions map[string]Session
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[string]Session),
	}
}

//	将客户端（连接）注册至 Hub 中
func (p *Hub) Register(serialNo string, session Session) error {
	if client := p.Find(serialNo); client != nil {
		if err := p.UnRegister(client); err != nil {
			return err
		}
	}

	p.rmu.Lock()
	defer p.rmu.Unlock()

	p.sessions[serialNo] = session

	return nil
}

func (p *Hub) UnRegister(client Session) error {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	//	Close Old Connection
	client.Close()
	delete(p.sessions, client.SerialNo())

	return nil
}

func (p *Hub) Find(serialNo string) Session {
	p.rmu.RLock()
	defer p.rmu.RUnlock()

	return p.sessions[serialNo]
}

func (p *Hub) Search(serialsNo ...string) BroadcastTarget {
	var result BroadcastTarget

	//	当未指定广播设备序列时，则使用完全广播模式
	if len(serialsNo) <= 0 {
		for _, session := range p.sessions {
			result = append(result, session)
		}
	}

	//	当指定广播设备序列时，则使用局部广播模式
	for _, serialNo := range serialsNo {
		if session, ok := p.sessions[serialNo]; ok {
			result = append(result, session)
		}
	}

	return result
}

func (p *Hub) Start(session Session) {
	//	开始消息循环
	session.WriterServ()
	session.ReaderServ()

	log.Printf("开始监听客户端 %s 状态\n", session.SerialNo())

	//	如果是收到断开的消息，则删除当前客户端
	if <-session.Dead() {
		p.rmu.Lock()
		defer p.rmu.Unlock()
		delete(p.sessions, session.SerialNo())
		log.Printf("检测到客户端 %s 离线，已清除服务器中的Session信息\n", session.SerialNo())
	}
}
