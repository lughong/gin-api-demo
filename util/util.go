package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/gin-gonic/gin"
)

// GetReqID 获取设置的X-Request-Id值
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

// EncryptMD5 加密纯文本，返回加密后的字符串
func EncryptMD5(source string) (string, error) {
	md := md5.New()

	if _, err := md.Write([]byte(source)); err != nil {
		return "", err
	}

	return hex.EncodeToString(md.Sum(nil)), nil
}

// PathExists 判断文件夹或文件是否存在
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

// CreateFile 创建文件
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

// CreateDir 创建文件夹
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
