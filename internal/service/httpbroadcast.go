package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SendRequest struct {
	SerialNo []string        `json:"serial_no" form:"serial_no"`
	Token    string          `json:"token" form:"token"`
	Message  json.RawMessage `json:"message" form:"message"`
}

func (p *WebService) broadcast(c *gin.Context) {
	var request SendRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"err_no": 1, "msg": err.Error()})
		return
	}

	if len(request.SerialNo) <= 0 {
		if p.token == "" || request.Token != p.token {
			c.JSON(http.StatusOK, gin.H{"err_no": 1, "msg": "发送失败，无效的Token"})
			return
		}
	}

	devices := p.hub.Search(request.SerialNo...)

	//	广播发送消息
	devices.Broadcast(request.Message)

	c.JSON(http.StatusOK, gin.H{"err_no": 0, "data": devices.GetSerialsNo(), "msg": "Queued"})
}
