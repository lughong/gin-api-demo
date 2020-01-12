package http_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/lughong/gin-api-demo/api/user/handler/http"
	myMock "github.com/lughong/gin-api-demo/api/user/mock"
	"github.com/lughong/gin-api-demo/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetByUsername(t *testing.T) {
	username := "zhangsan"
	moUser := model.NewUser(1, username, "", 19)

	mockLogic := new(myMock.Logic)
	mockLogic.On("GetByUsername", mock.Anything, username).Return(moUser, nil)

	userHandler := handler.UserHandler{
		UserLogic: mockLogic,
	}

	var data = struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: "",
	}
	body, err := json.Marshal(data)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = httptest.NewRequest("GET", "/user", bytes.NewReader(body))

	userHandler.GetByUsername(c)

	type expected struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	var exp = expected{
		Code: 0,
		Msg:  "OK",
		Data: map[string]interface{}{
			"username": moUser.GetUsername(),
			"age":      float64(moUser.GetAge()),
		},
	}

	result := resp.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var actual expected
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)
	assert.Equal(t, exp, actual)
	mockLogic.AssertExpectations(t)
}
