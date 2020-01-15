package auth

import "golang.org/x/crypto/bcrypt"

// Encrypt 加密纯文本，返回加密后的字符串
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare 加密文本和纯文本对比，如果一样则返回nil
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
