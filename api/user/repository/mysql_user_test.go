package repository_test

import (
	"context"
	"testing"

	"github.com/lughong/gin-api-demo/model"

	"github.com/lughong/gin-api-demo/api/user/repository"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
)

func TestGetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error %s was not expected when opening a stub database connection", err)
	}

	mockUser := model.NewUser(1, "zhangsan", "", 18)
	rows := sqlmock.NewRows([]string{"id", "username", "password", "age"}).
			AddRow(mockUser.GetID(), mockUser.GetUsername(), mockUser.GetPassword(), mockUser.GetAge())

	query := "SELECT id, username, password, age FROM user WHERE username=?"
	mock.ExpectPrepare(query).ExpectQuery().WillReturnRows(rows)

	userRepo := repository.NewMysqlUserRepository(db)
	anUser, err := userRepo.GetByUsername(context.TODO(), mockUser.GetUsername())
	assert.NoError(t, err)
	assert.NotNil(t, anUser)
}
