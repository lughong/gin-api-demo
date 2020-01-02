package repository

import (
	"context"
	"database/sql"

	"github.com/lughong/gin-api-demo/api/user"
	"github.com/lughong/gin-api-demo/model"
)

const (
	// SQL语句常量
	selectUserSQL       = "SELECT id, username, password, age FROM user"
	selectByUsernameSQL = " WHERE username=?"
	selectUserDetail    = selectUserSQL + selectByUsernameSQL
)

// mysqlUserRepository
type mysqlUserRepository struct {
	DB *sql.DB
}

// NewMysqlUserRepository 返回仓库接口
func NewMysqlUserRepository(db *sql.DB) user.Repository {
	return &mysqlUserRepository{
		DB: db,
	}
}

// GetByUsername 获取一条user数据
func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	stmt, err := m.DB.PrepareContext(ctx, selectUserDetail)
	if err != nil {
		return nil, err
	}

	var (
		id       int
		uname 	 string
		password string
		age      int
	)
	if err := stmt.QueryRow(username).Scan(&id, &uname, &password, &age); err != nil {
		return nil, err
	}

	// 把查找结果解析到User结构体
	anUser := model.NewUser(id, uname, password, age)
	return anUser, nil
}
