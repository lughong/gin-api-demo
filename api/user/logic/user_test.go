package logic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lughong/gin-api-demo/api/user/logic"
	myMock "github.com/lughong/gin-api-demo/api/user/mock"
	"github.com/lughong/gin-api-demo/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetByUsername(t *testing.T) {
	username := "zhangsan"
	mockRepo := &myMock.Repository{}

	t.Run("success", func(t *testing.T) {
		moUser := model.NewUser(1, username, "", 18)
		mockRepo.On("GetByUsername", mock.Anything, username).Return(moUser, nil).Once()

		userLogic := logic.NewUserLogic(mockRepo, time.Duration(5)*time.Second)
		user, err := userLogic.GetByUsername(context.TODO(), username)

		assert.NoError(t, err)
		assert.NotNil(t, user)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRepo.On("GetByUsername", mock.Anything, username).Return(nil, fmt.Errorf("some error")).Once()

		userLogic := logic.NewUserLogic(mockRepo, time.Duration(5)*time.Second)
		user, err := userLogic.GetByUsername(context.TODO(), username)

		assert.Error(t, err)
		assert.Nil(t, user)

		mockRepo.AssertExpectations(t)
	})
}
