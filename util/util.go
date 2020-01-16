package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
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

// Exists 判断文件夹或文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// CreateFile 创建文件
func CreateFile(fileName string) (bool, error) {
	if exists := Exists(fileName); exists {
		return true, nil
	}

	if _, err := os.Create(fileName); err != nil {
		return false, err
	}

	return true, nil
}

// CreateDir 创建文件夹
func CreateDir(path string) (bool, error) {
	if exists := Exists(path); exists {
		return true, nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return false, err
	}

	return true, nil
}

// SendMail 发送邮件
func SendMail(mailTo []string, subject string, body string) error {
	// 是否已经开启发送开关
	if !viper.GetBool("mail.sendOpen") {
		return nil
	}

	m := gomail.NewMessage()

	// 设置发件人
	m.SetHeader("From", viper.GetString("mail.from"))

	// 设置发送给多个用户
	m.SetHeader("To", mailTo...)

	// 设置邮件主题
	m.SetHeader("Subject", subject)

	// 设置邮件正文
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		viper.GetString("mail.host"),
		viper.GetInt("mail.port"),
		viper.GetString("mail.user"),
		viper.GetString("mail.password"),
	)

	return d.DialAndSend(m)
}
