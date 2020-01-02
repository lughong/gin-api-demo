package user

import (
	"context"

	"github.com/lughong/gin-api-demo/model"
)

// Logic 定义user的逻辑接口
type Logic interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}
