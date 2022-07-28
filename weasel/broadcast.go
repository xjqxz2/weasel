package weasel

type BroadcastTarget []Client

func (p BroadcastTarget) Broadcast(message []byte) {
	for _, client := range p {
		go func(client Client) {
			client.Write(message)
		}(client)
	}
}
