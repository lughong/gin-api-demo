package repository

import (
	"context"
	"database/sql"

	"github.com/lughong/gin-api-demo/api/user"
	"github.com/lughong/gin-api-demo/model"
)

const (
	// SQL语句常量
	selectUserSQL = "SELECT id, username, password, age FROM user"

	selectByUsernameSQL        = " WHERE username=?"
	selectUserDetailByUsername = selectUserSQL + selectByUsernameSQL

	selectByUserIDSQL        = " WHERE id=?"
	selectUserDetailByUserID = selectUserSQL + selectByUserIDSQL
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

// GetByUserID 获取一条user数据
func (m *mysqlUserRepository) GetByUserID(ctx context.Context, id int) (*model.User, error) {
	return m.getUser(ctx, selectUserDetailByUserID, id)
}

// GetByUsername 获取一条user数据
func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return m.getUser(ctx, selectUserDetailByUsername, username)
}

func (m *mysqlUserRepository) getUser(ctx context.Context, sql string, args ...interface{}) (*model.User, error) {
	stmt, err := m.DB.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	var (
		id       int
		username string
		password string
		age      int
	)
	if err := stmt.QueryRow(args...).Scan(&id, &username, &password, &age); err != nil {
		return nil, err
	}

	// 把查找结果解析到User结构体
	anUser := model.NewUser(id, username, password, age)
	return anUser, nil
}
