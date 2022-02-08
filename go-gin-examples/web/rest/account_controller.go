package rest

import (
	"cathub.me/go-gin-examples/domain/user"
	"cathub.me/go-gin-examples/pkg/errors"
	"cathub.me/go-gin-examples/pkg/response"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	userAuthenticator user.UserAuthenticator
	userService       user.UserService
}

func RegisterAccountController(gin *gin.RouterGroup) {
	controller := AccountController{userService: user.GetUserService(), userAuthenticator: user.GetUserAuthorizeService()}

	gin.POST("/register", controller.register)
	gin.POST("/authenticate", controller.authorize)
	gin.GET("/activation/:code", controller.activation)
	gin.POST("/resendActivateEmail", controller.resendActivateEmail)
}

// 注册用户
func (u *AccountController) register(ctx *gin.Context) {
	var registerUserDTO user.RegisterUserDTO
	if err := ctx.ShouldBindJSON(&registerUserDTO); err != nil {
		response.Fail(ctx, err)
		return
	}

	newUser, err := u.userService.Register(ctx, registerUserDTO)
	if err != nil {
		response.Fail(ctx, err)
		return
	}

	response.Ok(ctx, response.ResponseEntityBuilder().Body(newUser).Build())
}

// 认证, 成功返回 token
func (u *AccountController) authorize(ctx *gin.Context) {
	loginVM := user.LoginVM{}
	if err := ctx.ShouldBindJSON(&loginVM); err != nil {
		response.Fail(ctx, err)
		return
	}

	jwtToken, err := u.userAuthenticator.Authorize(ctx, loginVM)
	if err != nil {
		response.Fail(ctx, err)
		ctx.Abort()
		return
	}

	response.Ok(ctx, response.ResponseEntityBuilder().Body(jwtToken).Header("Authorization", jwtToken.Token).Build())
}

func (u *AccountController) activation(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		response.Fail(ctx, errors.BadRequest.Descf("Path param 'code' cannot be empty"))
		return
	}

	err := u.userService.Activation(ctx, code)
	if err != nil {
		response.Fail(ctx, err)
		return
	}

	response.Ok(ctx, response.ResponseEntityBuilder().Ok("Activation is successful, please log in").Build())
}

func (u *AccountController) resendActivateEmail(ctx *gin.Context) {
	param := make(map[string]string, 1)
	if err := ctx.ShouldBindJSON(&param); err != nil {
		response.Fail(ctx, err)
		return
	}

	email := param["email"]
	if email == "" {
		response.Fail(ctx, errors.BadRequest.Descf("Body param 'email' cannot be empty"))
		return
	}

	if err := u.userService.ResendActivateEmail(ctx, email); err != nil {
		response.Fail(ctx, err)
		return
	}

	response.Ok(ctx, response.ResponseEntityBuilder().Ok("Resend successfully, please check your email").Build())
}
