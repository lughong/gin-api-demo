package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/lughong/gin-api-demo/model"
)

type UserRepository struct {
	mock.Mock
}

// GetByUsername provides a mock function with given fields: ctx, username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	args := r.Called(ctx, username)

	var r0 model.User
	if rf, ok := args.Get(0).(func(context.Context, string) model.User); ok {
		r0 = rf(ctx, username)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(model.User)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, model.User
func (r *UserRepository) Create(ctx context.Context, u model.User) (int64, error) {
	args := r.Called(ctx, u)

	var r0 int64
	if rf, ok := args.Get(0).(func(context.Context, model.User) int64); ok {
		r0 = rf(ctx, u)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(int64)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}
