package api

import (
	"net/http"

	"github.com/lughong/gin-api-demo/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Response 响应请求的结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, msg := errno.DecodeErr(err)

	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
