package service

import (
	"database/sql"

	"github.com/lughong/gin-api-demo/model"
	"github.com/lughong/gin-api-demo/pkg/errno"
	"github.com/sirupsen/logrus"
)

// GetUserDetail 根据用户名和密码获取用户详情
func GetUserDetail(username, password string) (model.User, error) {
	userModel := model.NewUser(func(u *model.User) {
		u.Username = username
		u.Password = password
	})

	user, err := userModel.Find()
	if err == sql.ErrNoRows {
		return user, errno.ErrUserNotFound
	}

	// 如果系统错误，记录日志
	if err != nil {
		err = errno.New(errno.ErrGetUserDetail, err).Addf("username=?, password=?", user.Username, user.Password)
		logrus.Errorf("Get an error. %s", err)
		return user, errno.InternalServerError
	}

	return user, nil
}
