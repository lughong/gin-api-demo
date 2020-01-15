package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	// ErrMissingHeader 设置请求头Authorization内容为空的错误提示
	ErrMissingHeader   = errors.New("The length of the Authorization header is zero. ")
	ErrNotJwtMapClaims = errors.New("The token Claims is not jwt.MapClaims. ")
)

// Context jwt令牌的内容
type Context struct {
	ID       float64
	Username string
}

// Sign 使用密钥签名加密token内容
func Sign(c *gin.Context, ctx Context, secret string) (tokenString string, err error) {
	// load the jwt secret from the server config if the secret isn't specified
	if secret == "" {
		secret = viper.GetString("server.jwtSecret")
	}

	var expire time.Duration
	if expire, err = time.ParseDuration(viper.GetString("token.timeout")); err != nil {
		expire = 3600 * time.Second
	}

	// 设置token内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       ctx.ID,
		"username": ctx.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(expire).Unix(),
	})

	// 使用密钥对token内容做签名加密
	tokenString, err = token.SignedString([]byte(secret))
	return
}

// secretFunc 验证密钥格式.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// 验证token的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}

// Parse 使用特定的密钥验证token是否有效
// 如果token有效，则返回token内容
func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	// 解析token
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil || !token.Valid {
		return ctx, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx, ErrNotJwtMapClaims
	}

	// 读取token内容
	if id, ok := claims["id"].(float64); ok {
		ctx.ID = id
	}
	if username, ok := claims["username"].(string); ok {
		ctx.Username = username
	}

	return ctx, nil
}

// ParseRequest 解析请求，获取token内容
func ParseRequest(c *gin.Context) (*Context, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return &Context{}, ErrMissingHeader
	}

	// 加载jwt加密密钥
	secret := viper.GetString("server.jwtSecret")

	var tokenString string
	// 解析header头获取token内容
	_, _ = fmt.Sscanf(header, "Bearer %s", &tokenString)
	return Parse(tokenString, secret)
}
