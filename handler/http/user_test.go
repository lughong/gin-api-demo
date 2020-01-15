package http_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_http "github.com/lughong/gin-api-demo/handler/http"
	myMock "github.com/lughong/gin-api-demo/mock"
	"github.com/lughong/gin-api-demo/model"
)

func TestGet(t *testing.T) {
	var moUser model.User
	err := faker.FakeData(&moUser)
	assert.NoError(t, err)

	mockLogic := new(myMock.UserLogic)
	mockLogic.On("GetByUsername", mock.Anything, moUser.Username).Return(moUser, nil)

	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = httptest.NewRequest("GET", "/user/"+moUser.Username, nil)
	c.Params = gin.Params{gin.Param{Key: "username", Value: moUser.Username}}

	userHandler := _http.UserHandler{
		UserLogic: mockLogic,
	}
	userHandler.GetByUsername(c)

	result := resp.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var actual struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, 0, actual.Code)
	mockLogic.AssertExpectations(t)
}
