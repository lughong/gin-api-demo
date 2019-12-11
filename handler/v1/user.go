package v1

import (
	. "github.com/lughong/gin-api-demo/handler"
	"github.com/lughong/gin-api-demo/pkg/auth"
	"github.com/lughong/gin-api-demo/pkg/errno"
	"github.com/lughong/gin-api-demo/service"

	"github.com/gin-gonic/gin"
)

// GetUser 获取用户详情
func GetUser(c *gin.Context) {
	var r CreateRequest

	// 绑定数据到结构体
	if err := c.BindJSON(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	// 检测必填项
	if err := r.Validate(); err != nil {
		//SendResponse(c, errno.ErrValidation, nil)
		//return
	}

	// 获取用户详情
	user, err := service.GetUserDetail(r.Username, r.Password)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}

	// 校验密码
	if err := auth.Compare(user.Password, r.Password); err != nil {
		//SendResponse(c, errno.ErrPasswordIncorrect, nil)
		//return
	}

	SendResponse(c, nil, user)
}
