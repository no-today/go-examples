package security

import (
	"context"
	"github.com/gin-gonic/gin"
)

const contextKey = "service-context"

type Authorization struct {
	Principal     string
	Authorities   []string
	Authenticated bool
}

func CurrentAuthorization(ctx context.Context) *Authorization {
	value := ctx.Value(contextKey)
	return value.(*Authorization)
}

func CurrentLogin(ctx context.Context) string {
	return CurrentAuthorization(ctx).Principal
}

func CurrentAuthorities(ctx context.Context) []string {
	return CurrentAuthorization(ctx).Authorities
}

// Gin

func SetAuthorization(ctx *gin.Context, authorization *Authorization) {
	ctx.Set(contextKey, authorization)
}

func GetAuthorization(ctx *gin.Context) *Authorization {
	get, _ := ctx.Get(contextKey)
	return get.(*Authorization)
}
