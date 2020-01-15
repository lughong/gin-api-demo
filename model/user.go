package model

import "context"

type (
	// User 结构体
	User struct {
		ID       int
		Username string
		Password string
		Age      int
	}

	// UserLogic 定义user的逻辑接口
	UserLogic interface {
		GetByUsername(ctx context.Context, username string) (User, error)
		Create(ctx context.Context, user User) (int64, error)
	}

	// UserRepository 定义user的仓库接口
	UserRepository interface {
		GetByUsername(ctx context.Context, username string) (User, error)
		Create(ctx context.Context, user User) (int64, error)
	}
)
