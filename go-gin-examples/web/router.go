package web

import (
	"cathub.me/go-gin-examples/middleware/security/jwt"
	"cathub.me/go-gin-examples/web/rest"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(jwt.Jwt())

	rest.RegisterAccountController(r.Group("/api/v1"))

	// 必须登陆
	group := r.Group("/api/v1")
	group.Use(jwt.RequireAuthenticated())
	rest.RegisterUserResources(group)

	return r
}
