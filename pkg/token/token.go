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
	// ErrMissingHeader means the Authorization header was empty.
	ErrMissingHeader = errors.New("The length of the Authorization header is zero. ")
)

// Context is the context of the JSON web token.
type Context struct {
	ID       string
	Email    string
	Username string
}

// Sign signs the context with the specified secret.
func Sign(c *gin.Context, ctx Context, secret string) (tokenString string, err error) {
	// load the jwt secret from the server config if the secret isn't specified
	if secret == "" {
		secret = viper.GetString("server.jwtSecret")
	}

	// the token context
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       ctx.ID,
		"email":    ctx.Email,
		"username": ctx.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Unix,
	})

	// sign the token with the specified secret.
	tokenString, err = token.SignedString(secret)
	return
}

// secretFunc validate the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// make sure the alg is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	}
}

// Parse validates the token with the specified secret,
// and return the context if the token was valid.
func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))
	if err != nil {
		return ctx, err

		// Read the token if it's valid.
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.ID = claims["id"].(string)
		ctx.Email = claims["email"].(string)
		ctx.Username = claims["username"].(string)
		return ctx, nil

		// Other errors.
	} else {
		return ctx, err
	}
}

// ParseRequest parse context with the gin request context.
func ParseRequest(c *gin.Context) (*Context, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return &Context{}, ErrMissingHeader
	}

	// load the jwt secret from server config.
	secret := viper.GetString("server.jwtSecret")

	var tokenString string
	// parse the header to get the token part.
	_, _ = fmt.Sscanf(header, "Bearer %s", &tokenString)
	return Parse(tokenString, secret)
}
