package user

import (
	"context"

	"github.com/lughong/gin-api-demo/model"
)

// Repository 定义user的仓库接口
type Repository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}
