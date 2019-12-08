package entity

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lughong/gin-api-demo/pkg/errno"
)

type CreateRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"pwd" form:"pwd" validate:"required"`
}

type CreateResponse struct {
	Url string `json:"url"`
}

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
