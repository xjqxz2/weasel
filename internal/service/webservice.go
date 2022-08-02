package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"peon.top/weasel/internal/weasel"
)

type WebService struct {
	ip        string
	port      int
	debugMode bool
	engine    *gin.Engine
	hub       *weasel.Hub
}

func New(ip string, port int) *WebService {
	service := &WebService{
		engine: gin.Default(),
		ip:     ip,
		port:   port,
		hub:    weasel.NewHub(),
	}

	//	加载跨域中间件
	service.engine.Use(cors())

	//	加载路由
	service.loadRoutes()

	return service
}

func (p *WebService) loadRoutes() {
	p.engine.GET("/dev/conn", p.upgradeWebsocket)
	p.engine.POST("/msg/broadcast", p.broadcast)
	p.engine.Static("/example", "/home/xjqxz/Workplace/weasel/wssdk")
}

func (p *WebService) Listen() error {
	return p.engine.Run(fmt.Sprintf("%s:%d", p.ip, p.port))
}
