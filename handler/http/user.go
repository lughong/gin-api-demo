package http

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/lughong/gin-api-demo/global/errno"
	"github.com/lughong/gin-api-demo/global/token"
	"github.com/lughong/gin-api-demo/handler"
	"github.com/lughong/gin-api-demo/model"
)

type UserHandler struct {
	UserLogic model.UserLogic
}

func NewUserHandler(g *gin.Engine, userLogic model.UserLogic) {
	h := &UserHandler{
		UserLogic: userLogic,
	}

	g.GET("/user/:username", h.GetByUsername)
	g.POST("/user", h.Create)
	g.POST("/login", h.Login)
}

// GetByUsername 获取用户详情
// @Summary Get user info of the database
// @Description Get user info
// @Tags user
// @Access json
// @Produce json
// @Success 200 {object} handler.Response "{"code":0,"msg":"OK","data":{"username":"zhangsan","age":18}}"
// @Router /user/:username [get]
func (u *UserHandler) GetByUsername(c *gin.Context) {
	// 获取用户详情
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	username := c.Param("username")
	if username == "" {
		handler.SendResponse(c, errno.ErrInvalidArgs, nil)
		return
	}

	anUser, err := u.UserLogic.GetByUsername(ctx, username)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	data := map[string]interface{}{
		"username": anUser.Username,
		"age":      anUser.Age,
	}
	handler.SendResponse(c, nil, data)
}

// Create 创建用户
// @Summary Create user to save the database
// @Description Create user info
// @Tags user
// @Access json
// @Produce json
// @Param username body handler.CreateRequest true "Get user info for the username"
// @Param password body handler.CreateRequest true "Get user info for the password"
// @Success 200 {object} handler.Response "{"code":0,"msg":"OK","data":{"username":"zhangsan"}}"
// @Router /user [post]
func (u *UserHandler) Create(c *gin.Context) {
	var r handler.CreateRequest

	// 绑定数据到结构体
	if err := c.BindJSON(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// 检测必填项
	if err := r.Validate(); err != nil {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	if err := r.Encrypt(); err != nil {
		handler.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// 获取用户详情
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	anUser, err := u.UserLogic.GetByUsername(ctx, r.Username)
	if err != nil && !errno.IsErrUserNotFound(err) {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}
	if anUser.Username != "" {
		handler.SendResponse(c, errno.ErrUserAlreadyExists, nil)
		return
	}

	user := model.User{
		Username: r.Username,
		Password: r.Password,
		Age:      r.Age,
	}
	_, err = u.UserLogic.Create(ctx, user)
	if err != nil {
		handler.SendResponse(c, errno.ErrCreateUser, nil)
		return
	}

	data := map[string]interface{}{
		"username": r.Username,
	}
	handler.SendResponse(c, nil, data)
}

// Login 用户登录验证
// @Summary user Login of the database
// @Description user Login
// @Tags user
// @Access json
// @Produce json
// @Param username body handler.CreateRequest true "Get user info for the username"
// @Param password body handler.CreateRequest true "Get user info for the password"
// @Success 200 {object} handler.Response "{"code":0,"msg":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1Nzg5OTU0NTMsImlkIjozLCJuYmYiOjE1Nzg5OTU0NTMsInVzZXJuYW1lIjoibGlzaSJ9.agmaafda4LwOqkqDbIkpV9AHkdaoFVHhOMkasu_qCTM"}}"
// @Router /login [post]
func (u *UserHandler) Login(c *gin.Context) {
	var r handler.CreateRequest

	// 绑定数据到结构体
	if err := c.BindJSON(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// 检测必填项
	if err := r.Validate(); err != nil {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	if err := r.Encrypt(); err != nil {
		handler.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	// 获取用户详情
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	anUser, err := u.UserLogic.GetByUsername(ctx, r.Username)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	if anUser.Password != r.Password {
		handler.SendResponse(c, errno.ErrPasswordInvalid, nil)
		return
	}

	cc := token.Context{
		ID:       float64(anUser.ID),
		Username: anUser.Username,
	}

	sign, err := token.Sign(c, cc, "")
	if err != nil {
		handler.SendResponse(c, errno.ErrToken, nil)
		return
	}

	data := map[string]interface{}{
		"token": "Bearer " + sign,
	}
	handler.SendResponse(c, nil, data)
}
