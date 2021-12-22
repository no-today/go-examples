package jwt

import (
	"cathub.me/go-web-examples/middleware/security"
	"cathub.me/go-web-examples/pkg/errors"
	"cathub.me/go-web-examples/pkg/response"
	"cathub.me/go-web-examples/pkg/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

// Claims 认证信息
type Claims struct {
	Principal   string
	Authorities []string
	jwt.StandardClaims
}

type JWTToken struct {
	Token string `json:"token,omitempty"`
}

func GenerateToken(claims Claims) (*JWTToken, error) {
	return generateToken(claims, setting.Jwt.Base64Secret, setting.Jwt.TokenValidityInSeconds)
}

func generateToken(claims Claims, jwtSecret string, tokenValidityInSeconds int64) (*JWTToken, error) {
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Principal:   claims.Principal,
		Authorities: claims.Authorities,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(tokenValidityInSeconds) * time.Second).Unix(),
		},
	}).SignedString([]byte(jwtSecret))

	return &JWTToken{Token: fmt.Sprintf("Bearer %s", token)}, nil
}

func ParseToken(token string) (*jwt.Token, error) {
	return parseToken(token, setting.Jwt.Base64Secret)
}

func parseToken(token string, jwtSecret string) (*jwt.Token, error) {
	bearerAndToken := strings.Split(token, "Bearer ")
	if len(bearerAndToken) < 2 {
		return nil, errors.Unauthorized
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(bearerAndToken[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.Unauthorized.Err(err)
		}
		validationError, ok := err.(*jwt.ValidationError)
		if ok {
			if validationError.Errors == jwt.ValidationErrorExpired {
				return nil, errors.TokenExpired.Err(validationError)
			}
		}
		return nil, errors.Unauthorized.Err(err)
	}

	if !tkn.Valid {
		return nil, errors.Unauthorized
	}

	return tkn, nil
}

const Authorization = "Authorization"

func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(Authorization)
		if token == "" {
			token = ctx.Query(Authorization)
		}

		_jwt, err := ParseToken(token)
		if err != nil {
			// 匿名用户
			security.SetAuthorization(ctx, &security.Authorization{
				Authorities:   []string{security.ANONYMOUS},
				Authenticated: false,
			})
		} else {
			// 认证用户
			claims := _jwt.Claims.(*Claims)
			security.SetAuthorization(ctx, &security.Authorization{
				Principal:     claims.Principal,
				Authorities:   claims.Authorities,
				Authenticated: true,
			})
		}
	}
}

func RequireAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !security.GetAuthorization(ctx).Authenticated {
			response.Fail(ctx, errors.Unauthorized)
			ctx.Abort()
		}
	}
}
