package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"peon.top/weasel/weasel"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (p *WebService) upgradeWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	var (
		serialNo   = c.Query("serial_no")
		serialName = c.Query("serial_name")
	)

	//	创建一个客户端
	session := weasel.NewWSSession(conn, serialNo, serialName)

	//	释放客户端连接资源
	defer session.Close()
	log.Printf("客户端 %s 已连接到服务器，已为其开启数据收发服务\n", c.Query("serial_no"))

	if err := p.hub.Register(serialNo, session); err != nil {
		return
	}

	p.hub.Start(session)
}
