package handler

import (
	"net/http"

	"github.com/lughong/gin-api-demo/app/pkg/auth"
	"github.com/lughong/gin-api-demo/app/pkg/errno"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

// CreateRequest 接收请求结构体信息
type CreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Validate 检测必填项信息
func (r *CreateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Encrypt 加密密码
func (r *CreateRequest) Encrypt() (err error) {
	r.Password, err = auth.Encrypt(r.Password)
	return
}

// CreateResponse 结构体设置了响应内容
type CreateResponse struct {
	Username string `json:"username"`
}

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
