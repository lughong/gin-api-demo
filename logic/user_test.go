package logic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lughong/gin-api-demo/logic"
	myMock "github.com/lughong/gin-api-demo/mock"
	"github.com/lughong/gin-api-demo/model"
)

func TestGetByUsername(t *testing.T) {
	mockRepo := &myMock.UserRepository{}

	t.Run("success", func(t *testing.T) {
		var moUser model.User
		err := faker.FakeData(&moUser)
		assert.NoError(t, err)

		mockRepo.On("GetByUsername", mock.Anything, moUser.Username).Return(moUser, nil).Once()

		userLogic := logic.NewUserLogic(mockRepo, time.Duration(5)*time.Second)
		user, err := userLogic.GetByUsername(context.TODO(), moUser.Username)

		assert.NoError(t, err)
		assert.NotNil(t, user)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		var expected model.User
		mockRepo.On("GetByUsername", mock.Anything, expected.Username).Return(expected, fmt.Errorf("some error")).Once()

		userLogic := logic.NewUserLogic(mockRepo, time.Duration(5)*time.Second)
		user, err := userLogic.GetByUsername(context.TODO(), expected.Username)

		assert.Error(t, err)
		assert.Equal(t, user, expected)
		mockRepo.AssertExpectations(t)
	})
}
