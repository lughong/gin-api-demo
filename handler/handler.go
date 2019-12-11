package handler

import (
	"net/http"

	"github.com/lughong/gin-api-demo/pkg/auth"
	"github.com/lughong/gin-api-demo/pkg/errno"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

type CreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Validate 检测必填项
func (r *CreateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Encrypt 加密密码
func (r *CreateRequest) Encrypt() (err error) {
	r.Password, err = auth.Encrypt(r.Password)
	return
}

type CreateResponse struct {
	Username string `json:"username"`
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
