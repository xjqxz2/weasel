package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type KickRequest struct {
	SerialNo []string `json:"serial_no" form:"serial_no"`
}

func (p *WebService) kick(c *gin.Context) {
	var request KickRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"err_no": 1, "msg": err.Error()})
		return
	}

	if len(request.SerialNo) <= 0 {
		c.JSON(http.StatusOK, gin.H{"err_no": 1, "msg": "请指定一个客户端踢下线"})
		return
	}

	devices := p.hub.Search(request.SerialNo...)
	devices.Kick()

	c.JSON(http.StatusOK, gin.H{"err_no": 1, "msg": "请指定一个客户端踢下线"})
}
