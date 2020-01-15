package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"

	"github.com/lughong/gin-api-demo/repository"
)

func TestGetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error %s was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "username", "password", "age"}).
		AddRow(1, "zhangsan", "", 18).
		AddRow(2, "lisi", "", 19)

	query := "SELECT id, username, password, age FROM user WHERE username=?"
	mock.ExpectPrepare(query).ExpectQuery().WillReturnRows(rows)

	username := "zhangsan"
	userRepo := repository.NewMysqlUserRepository(db)
	anUser, err := userRepo.GetByUsername(context.TODO(), username)
	assert.NoError(t, err)
	if assert.NotNil(t, anUser) {
		assert.Equal(t, username, anUser.Username)
	}
}
