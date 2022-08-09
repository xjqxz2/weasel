package weasel

import (
	"fmt"
	"log"
	"net/http"
)

type Event interface {
	Fire(packet *EventPacket)
}

type NoNotifyEvent struct{}

func NewNoNotifyEvent() *NoNotifyEvent {
	return new(NoNotifyEvent)
}

func (p *NoNotifyEvent) Fire(packet *EventPacket) {}

type EventPacket struct {
	PackType int
	DeviceId string
}

type RemoteEvent struct {
	Host string
}

func NewRemoteEvent(host string) *RemoteEvent {
	return &RemoteEvent{Host: host}
}

func (p *RemoteEvent) Fire(packet *EventPacket) {
	go func() {
		client := &http.Client{}
		_, err := client.Get(fmt.Sprintf("%s?id=%s&type=%d", p.Host, packet.DeviceId, packet.PackType))

		if err != nil {
			log.Printf("通知失败:%s\n", err.Error())
			return
		}

		log.Printf("消息通知成功:%s -> PacketType(%d)\n", packet.DeviceId, packet.PackType)
	}()
}
