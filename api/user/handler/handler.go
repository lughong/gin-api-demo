package handler

import (
	"github.com/lughong/gin-api-demo/pkg/auth"

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
