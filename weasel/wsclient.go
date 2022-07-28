package weasel

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
)

//	The Websocket Client struct
//	Impl a Client interface
type WSClient struct {
	NetworkClient

	//	a websocket connection
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWSClient(conn *websocket.Conn, serialNo, serialName string) *WSClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &WSClient{
		NetworkClient: NewNetworkClient(serialNo, serialName),
		conn:          conn,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (p *WSClient) Write(b []byte) {
	p.MsgWriter <- b
}

func (p *WSClient) Receive() <-chan []byte {
	return p.MsgReader
}

func (p *WSClient) WriterServ() {
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

func (p *WSClient) ReaderServ() {
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
			default:
				log.Printf("接收到 %s 的信息，内容为 %s\n", p.serialNo, string(message))

				p.MsgReader <- message
			}
		}

	}
}

func (p *WSClient) Close() {
	if p.conn != nil {
		p.cancel()
		_ = p.conn.Close()
	}
}
