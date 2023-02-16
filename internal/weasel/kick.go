package weasel

func (p BroadcastTarget) Kick() {
	for _, session := range p {
		go func(session Session) {
			//	向客户端发送一条下线包，用于关闭定时器
			session.Write([]byte(PACK_KICK_OFFLINE))

			//	关闭服务端
			session.Close()
		}(session)
	}
}
