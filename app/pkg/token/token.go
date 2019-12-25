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
	ErrMissingHeader = errors.New("The length of the Authorization header is zero. ")
)

// Context jwt令牌的内容
type Context struct {
	ID       string
	Email    string
	Username string
}

// Sign 使用密钥签名加密token内容
func Sign(c *gin.Context, ctx Context, secret string) (tokenString string, err error) {
	// load the jwt secret from the server config if the secret isn't specified
	if secret == "" {
		secret = viper.GetString("server.jwtSecret")
	}

	// 设置token内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       ctx.ID,
		"email":    ctx.Email,
		"username": ctx.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Unix,
	})

	// 使用密钥对token内容做签名加密
	tokenString, err = token.SignedString(secret)
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
	if err != nil {
		return ctx, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 读取token内容
		ctx.ID = claims["id"].(string)
		ctx.Email = claims["email"].(string)
		ctx.Username = claims["username"].(string)
		return ctx, nil
	} else {
		return ctx, err
	}
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
