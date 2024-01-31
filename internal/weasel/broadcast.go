package weasel

import "fmt"

type Broadcasts []Session

func (p *Broadcasts) Broadcast(message []byte) {
	//	如果未找到客户端则直接返回，不再进行没有意义的广播
	if len(*p) <= 0 {
		return
	}

	go func() {
		for _, s := range *p {
			s.Write(message)
			fmt.Printf("准备下发[%s]消息到[%s]设备", string(message), s.SerialNo())
		}
	}()
}

func (p *Broadcasts) GetSerialsNo() (result []string) {
	if len(*p) <= 0 {
		return []string{}
	}

	for _, session := range *p {
		result = append(result, session.SerialNo())
	}

	return result
}
