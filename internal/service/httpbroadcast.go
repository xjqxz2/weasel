package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SendRequest struct {
	SerialNo []string        `json:"serial_no" form:"serial_no"`
	Message  json.RawMessage `json:"message" form:"message"`
}

func (p *WebService) broadcast(c *gin.Context) {
	var request SendRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, gin.H{"err_no": 1, "msg": err.Error()})
		return
	}

	devices := p.hub.Search(request.SerialNo...)

	//	广播发送消息
	devices.Broadcast(request.Message)

	c.JSON(http.StatusOK, gin.H{"err_no": 0, "data": devices.GetSerialsNo(), "msg": "Queued"})
}
