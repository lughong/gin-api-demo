package mock

import (
	"context"

	"github.com/lughong/gin-api-demo/model"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

// GetByUsername provides a mock function with given fields: ctx, username
func (r *Repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := r.Called(ctx, username)

	var r0 *model.User
	if rf, ok := args.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, username)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*model.User)
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

// GetByUserID provides a mock function with given fields: ctx, username
func (r *Repository) GetByUserID(ctx context.Context, id int) (*model.User, error) {
	args := r.Called(ctx, id)

	var r0 *model.User
	if rf, ok := args.Get(0).(func(context.Context, int) *model.User); ok {
		r0 = rf(ctx, id)
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}
