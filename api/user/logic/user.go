package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/lughong/gin-api-demo/api/user"
	"github.com/lughong/gin-api-demo/model"
	"github.com/lughong/gin-api-demo/pkg/errno"

	"github.com/sirupsen/logrus"
)

type userLogic struct {
	userRepo       user.Repository
	contextTimeout time.Duration
}

func NewUserLogic(repo user.Repository, timeout time.Duration) user.Logic {
	return &userLogic{
		userRepo:       repo,
		contextTimeout: timeout,
	}
}

// GetByUsername 根据用户名获取用户详情
func (u *userLogic) GetByUsername(c context.Context, username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	anUser, err := u.userRepo.GetByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return nil, errno.ErrUserNotFound
	}

	// 如果系统错误，记录日志
	if err != nil {
		err = errno.New(errno.ErrGetUserDetail, err).Addf("username=?", username)
		logrus.Errorf("Get an error. %s", err)
		return nil, errno.InternalServerError
	}

	return anUser, nil
}
