package util

import (
	"os"

	"github.com/gin-gonic/gin"
)

func GetReqID(c *gin.Context) string {
	value, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}

	if requestID, ok := value.(string); ok {
		return requestID
	}

	return ""
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CreateFile(fileName string) (bool, error) {
	ok, err := PathExists(fileName)
	if err != nil {
		return false, err
	}

	if !ok {
		_, err := os.Create(fileName)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func CreateDir(path string) (bool, error) {
	ok, err := PathExists(path)
	if err != nil {
		return false, err
	}

	if !ok {
		if err := os.MkdirAll(path, 0755); err != nil {
			return false, err
		}
	}

	return true, nil
}
