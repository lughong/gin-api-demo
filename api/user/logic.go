package user

import (
	"context"

	"github.com/lughong/gin-api-demo/model"
)

// Logic 定义user的逻辑接口
type Logic interface {
	GetByUserID(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}
