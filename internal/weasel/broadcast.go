package weasel

type BroadcastTarget []Session

func (p BroadcastTarget) Broadcast(message []byte) {
	//	如果未找到客户端则直接返回，不再进行没有意义的广播
	if len(p) <= 0 {
		return
	}

	for _, client := range p {
		go func(s Session) {
			s.Write(message)
		}(client)
	}
}

func (p BroadcastTarget) GetSerialsNo() (result []string) {
	if len(p) <= 0 {
		return []string{}
	}

	for _, session := range p {
		result = append(result, session.SerialNo())
	}

	return result
}
