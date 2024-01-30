package weasel

import "github.com/gin-gonic/gin"

type RequestInfo struct {
	SerialNo   string `form:"serial_no" query:"serial_no" json:"serial_no"`
	SerialName string `form:"serial_name" query:"serial_name" json:"serial_name"`
	Code       string `form:"code" query:"code" json:"code"`
	Token      string `form:"token" query:"token" json:"token"`
	Path       string `form:"code" query:"path" json:"path"`
}

func (p *RequestInfo) Bind(c *gin.Context) error {
	return c.ShouldBindQuery(c)
}
