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

	p.hub.Search(request.SerialNo...).Broadcast(request.Message)

	c.JSON(http.StatusOK, gin.H{"err_no": 0, "msg": "Queued"})
}
