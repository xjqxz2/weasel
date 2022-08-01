package weasel

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
)

//	The Websocket Client struct
//	Impl a Client interface
type WSSession struct {
	NetworkClient

	//	a websocket connection
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWSSession(conn *websocket.Conn, serialNo, serialName string) *WSSession {
	ctx, cancel := context.WithCancel(context.Background())

	return &WSSession{
		NetworkClient: NewNetworkClient(serialNo, serialName),
		conn:          conn,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (p *WSSession) Write(b []byte) {
	p.MsgWriter <- b
}

func (p *WSSession) Receive() <-chan []byte {
	return p.MsgReader
}

func (p *WSSession) WriterServ() {
	go func() {
		select {
		case <-p.ctx.Done():
			log.Printf("客户端 %s 接收到停止信号，释放 WriterServ \n", p.serialNo)
			return
		default:
			for writer := range p.MsgWriter {
				//	向客户端写入数据
				if err := p.conn.WriteMessage(websocket.TextMessage, writer); err != nil {
					log.Printf("向客户端 %s 发送消息失败 -> %s\n", p.serialNo, err.Error())
					p.Close()
					continue
				}
			}
		}

	}()
}

func (p *WSSession) ReaderServ() {
	for {
		select {
		case <-p.ctx.Done():
			log.Printf("客户端 %s 接收到停止信号，释放 ReaderServ\n", p.serialNo)
			return
		default:
			messageType, message, err := p.conn.ReadMessage()

			//	如果消息出错了，则关闭连接
			if err != nil {
				p.Close()
				log.Printf("客户端 %s 接收消息失败 -> %s\n", p.serialNo, err.Error())
				continue
			}

			//	判定消息类型，做出对应的处理
			switch messageType {
			case websocket.PingMessage:
				_ = p.conn.WriteMessage(websocket.PongMessage, nil)
				log.Printf("已回复 %s 的心跳检测包 \n", p.serialNo)
			default:
				log.Printf("接收到 %s 的信息，内容为 %s\n", p.serialNo, string(message))
				//	暂时不将收到的消息放入管道中 对该部分还未进行实现，因为暂时还用不到此处的代码
				//p.MsgReader <- message
			}
		}

	}
}

func (p *WSSession) Close() {
	if p.conn != nil {
		p.cancel()
		_ = p.conn.Close()
	}
}
