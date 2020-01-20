package http

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/lughong/gin-api-demo/global/errno"
	"github.com/lughong/gin-api-demo/global/token"
	"github.com/lughong/gin-api-demo/handler"
	"github.com/lughong/gin-api-demo/model"
)

// UserHandler
type UserHandler struct {
	UserLogic model.UserLogic
}

// NewUserHandler
func NewUserHandler(userLogic model.UserLogic) *UserHandler {
	return &UserHandler{
		UserLogic: userLogic,
	}
}

// @Summary 获取用户信息
// @Description 从数据库中获取用户信息
// @Tags user
// @Access json
// @Produce json
// @param username path string true "Username"
// @Success 200 {string} string "{"code":0,"msg":"OK","data":{"username":"zhangsan","age":18}}"
// @Router /user/{username} [get]
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

// @Summary 新增一个用户
// @Description 成功后，返回新增用户名称
// @Tags user
// @Access json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Param age body int true "Password"
// @Success 200 {string} string "{"code":0,"msg":"OK","data":{"username":"zhangsan"}}"
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

// @Summary 用户登录
// @Description 登录成功返回一个token，后面访问操作都需要带上这个token值作校验
// @Tags Login
// @Access json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {string} string "{"code":0,"msg":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1Nzg5OTU0NTMsImlkIjozLCJuYmYiOjE1Nzg5OTU0NTMsInVzZXJuYW1lIjoibGlzaSJ9.agmaafda4LwOqkqDbIkpV9AHkdaoFVHhOMkasu_qCTM"}}"
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
