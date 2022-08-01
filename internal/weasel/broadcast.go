package weasel

type BroadcastTarget []Session

func (p BroadcastTarget) Broadcast(message []byte) {
	for _, client := range p {
		go func(s Session) {
			s.Write(message)
		}(client)
	}
}
