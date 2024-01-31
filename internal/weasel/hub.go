package weasel

import (
	"log"
	"sync"
)

type Hub struct {
	rmu      sync.RWMutex
	sessions map[string][]Session
	keeper   Keeper
	ev       Event
}

func NewHub(messageKeeper Keeper, remoteEvent Event) *Hub {
	return &Hub{
		sessions: make(map[string][]Session),
		keeper:   messageKeeper,
		ev:       remoteEvent,
	}
}

// 将客户端（连接）注册至 Hub 中
func (p *Hub) Register(serialNo string, session Session, requestInfo *RequestInfo) error {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	session.SetKeeper(p.keeper)
	sessions := p.sessions[serialNo]
	sessions = append(sessions, session)
	p.sessions[serialNo] = sessions

	//	通知使用事件机制（发送 连接包）
	p.ev.Fire(&EventPacket{
		PackType:    1,
		DeviceId:    session.SerialNo(),
		RequestInfo: *requestInfo},
	)

	return nil
}

func (p *Hub) UnRegister(session Session) error {
	p.rmu.Lock()
	defer p.rmu.Unlock()

	//	Close Old Connection
	session.Close()
	delete(p.sessions, session.SerialNo())

	return nil
}

func (p *Hub) Find(serialNo string) []Session {
	p.rmu.RLock()
	defer p.rmu.RUnlock()

	return p.sessions[serialNo]
}

func (p *Hub) Search(serialsNo ...string) Broadcasts {
	var result Broadcasts

	switch {
	case len(serialsNo) <= 0:
		//	当未指定广播设备序列时，则使用完全广播模式
		for _, session := range p.sessions {
			result = append(result, session...)
		}
	default:
		//	当指定广播设备序列时，则使用局部广播模式
		for _, serialNo := range serialsNo {
			if session, ok := p.sessions[serialNo]; ok {
				result = append(result, session...)
			}
		}
	}

	return result
}

func (p *Hub) Start(session Session) {
	//	开始消息循环
	session.WriterServ()
	session.ReaderServ()

	//	下发最新的一条消息
	if message := p.keeper.Message(session.SerialNo()); message != nil {
		session.Write(message)
	}

	log.Printf("开始监听客户端 %s 状态\n", session.SerialNo())

	//	如果是收到断开的消息，则删除当前客户端
	if <-session.Dead() {
		p.rmu.Lock()
		defer p.rmu.Unlock()

		//	获取设备下的所有 Session 对象
		if sessions, ok := p.sessions[session.SerialNo()]; ok {

			//	若该 SerialNo 下的 Session 数量 > 0 则找到当前这个，将其释放掉
			if len(sessions) > 0 {
				sessions = p.cleanSession(sessions, session.SessionId())
				p.sessions[session.SerialNo()] = sessions
			}

			if len(sessions) <= 0 {
				p.keeper.Offline(session.SerialNo())
				delete(p.sessions, session.SerialNo())
				log.Printf("客户端 %s 已没有可用的设备释放内存空间\n", session.SerialNo())
			}

			//	通知使用事件机制（发送 连接包）
			p.ev.Fire(&EventPacket{PackType: 2, DeviceId: session.SerialNo()})
			log.Printf("检测到客户端 %s 离线，已清除服务器中的Session信息\n", session.SerialNo())
		}
	}
}

func (p *Hub) cleanSession(sessions []Session, sessionId string) (result []Session) {
	for _, session := range sessions {
		if session.SessionId() != sessionId {
			result = append(result, session)
		}
	}

	return result
}
