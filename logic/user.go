package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/lughong/gin-api-demo/global/errno"
	"github.com/lughong/gin-api-demo/model"
)

// userLogic
type userLogic struct {
	userRepo       model.UserRepository
	contextTimeout time.Duration
}

// NewUserLogic
func NewUserLogic(repo model.UserRepository, timeout time.Duration) model.UserLogic {
	return &userLogic{
		userRepo:       repo,
		contextTimeout: timeout,
	}
}

// GetByUsername 根据用户名获取用户详情
func (u *userLogic) GetByUsername(c context.Context, username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	var user model.User
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return user, errno.ErrUserNotFound
	}

	// 如果系统错误，记录日志
	if err != nil {
		logrus.Errorf("UserLogic GetByUsername. %s", err)
		return user, err
	}

	return user, nil
}

// CreateUser 根据用户名获取用户详情
func (u *userLogic) Create(c context.Context, user model.User) (int64, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	lastInsertId, err := u.userRepo.Create(ctx, user)
	if err != nil {
		logrus.Errorf("UserLogic Create. %s", err)
		return 0, err
	}

	return lastInsertId, nil
}
