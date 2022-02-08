package rest

import (
	"cathub.me/go-gin-examples/domain/user"
	"cathub.me/go-gin-examples/pkg/data"
	"cathub.me/go-gin-examples/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserResources struct {
	userService user.UserService
}

func RegisterUserResources(gin *gin.RouterGroup) {
	resources := UserResources{userService: user.GetUserService()}

	gin.GET("/users", resources.findAll)
}

func (u *UserResources) findAll(ctx *gin.Context) {
	req := data.PageRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		response.Fail(ctx, err)
		return
	}

	pageable, users, err := u.userService.FindAll(ctx, data.ToPageable(&req))
	if err != nil {
		response.Fail(ctx, err)
		return
	}

	response.Ok(ctx, response.ResponseEntityBuilder().Body(data.ToPageResp(pageable, users)).Build())
}
