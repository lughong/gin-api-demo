package repository

import (
	"context"
	"database/sql"

	"github.com/lughong/gin-api-demo/model"
)

const (
	// SQL语句常量
	selectUserSQL = "SELECT id, username, password, age FROM user"

	selectByUsernameSQL        = " WHERE username=?"
	selectUserDetailByUsername = selectUserSQL + selectByUsernameSQL

	insertUserSQL = "INSERT INTO user (username, password, age) VALUES (?, ?, ?)"
)

// mysqlUserRepository
type mysqlUserRepository struct {
	DB *sql.DB
}

// NewMysqlUserRepository 返回仓库接口
func NewMysqlUserRepository(db *sql.DB) model.UserRepository {
	return &mysqlUserRepository{
		DB: db,
	}
}

// GetByUsername 获取一条user数据
func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	return m.getUser(ctx, selectUserDetailByUsername, username)
}

func (m *mysqlUserRepository) getUser(ctx context.Context, sql string, args ...interface{}) (model.User, error) {
	var user model.User

	stmt, err := m.DB.PrepareContext(ctx, sql)
	if err != nil {
		return user, err
	}

	if err := stmt.QueryRow(args...).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Age,
	); err != nil {
		return user, err
	}

	return user, nil
}

func (m *mysqlUserRepository) Create(ctx context.Context, user model.User) (int64, error) {
	stmt, err := m.DB.PrepareContext(ctx, insertUserSQL)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(
		user.Username,
		user.Password,
		user.Age,
	)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}
