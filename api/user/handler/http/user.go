package http

import (
	"context"

	"github.com/lughong/gin-api-demo/api/user"
	. "github.com/lughong/gin-api-demo/api/user/handler"
	"github.com/lughong/gin-api-demo/model"
	"github.com/lughong/gin-api-demo/pkg/auth"
	"github.com/lughong/gin-api-demo/pkg/errno"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserLogic user.Logic
}

func NewUserHandler(g *gin.Engine, userLogic user.Logic) {
	h := &UserHandler{
		UserLogic: userLogic,
	}

	g.GET("/user", h.GetByUsername)
}

// GetByUsername 获取用户详情
// @Summary Get user info of the database
// @Description Get user info
// @Tags user
// @Access json
// @Produce json
// @Param username body handler.CreateRequest true "Get user info for the username"
// @Param password body handler.CreateRequest true "Get user info for the password"
// @Success 200 {object} handler.Response "{"code":0,"msg":"OK","data":{"id":1,"username":"zhangsan","password":"","age":18}}"
// @Router /v1/user [get]
func (u *UserHandler) GetByUsername(c *gin.Context) {
	var r CreateRequest

	// 绑定数据到结构体
	if err := c.BindJSON(&r); err != nil {
		model.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// 检测必填项
	if err := r.Validate(); err != nil {
		//SendResponse(c, errno.ErrValidation, nil)
		//return
	}

	// 获取用户详情
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	anUser, err := u.UserLogic.GetByUsername(ctx, r.Username)
	if err != nil {
		model.SendResponse(c, err, nil)
		return
	}

	// 校验密码
	if err := auth.Compare(anUser.GetPassword(), r.Password); err != nil {
		//SendResponse(c, errno.ErrPasswordIncorrect, nil)
		//return
	}

	data := map[string]interface{}{
		"username": anUser.GetUsername(),
		"age":      anUser.GetAge(),
	}
	model.SendResponse(c, nil, data)
}
