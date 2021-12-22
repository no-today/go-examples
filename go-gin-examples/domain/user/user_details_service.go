package user

import (
	"cathub.me/go-web-examples/middleware/security/jwt"
	"cathub.me/go-web-examples/pkg/errors"
	"cathub.me/go-web-examples/pkg/utils/encrypt"
	"context"
	"regexp"
	"sync"
)

// _userAuthenticator Singleton
var _userAuthenticator UserAuthenticator
var _onceUserAuthenticator sync.Once

func GetUserAuthorizeService() UserAuthenticator {
	_onceUserAuthenticator.Do(func() {
		_userAuthenticator = &domainUserAuthenticator{userRepository: GetUserRepository()}
	})
	return _userAuthenticator
}

type UserDetails struct {
	Principal   string
	Credentials string
	Authorities []string
	Activated   bool
}

type UserAuthenticator interface {
	Authorize(ctx context.Context, vm LoginVM) (*jwt.JWTToken, error)
}

type domainUserAuthenticator struct {
	userRepository UserRepository
}

const ReEmail = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"

func (d domainUserAuthenticator) loadUserDetailByLogin(ctx context.Context, login string) (*UserDetails, error) {
	var user *User
	isEmailLogin, _ := regexp.MatchString(ReEmail, login)
	if isEmailLogin {
		user, _ = d.userRepository.FindByEmail(ctx, login)
	} else {
		user, _ = d.userRepository.FindByUsername(ctx, login)
	}

	if user == nil {
		return nil, errors.UserNotFound
	}

	return &UserDetails{
		Principal:   user.Username,
		Credentials: user.Password,
		Authorities: user.Roles,
		Activated:   user.Activated,
	}, nil
}

func (d domainUserAuthenticator) Authorize(ctx context.Context, vm LoginVM) (*jwt.JWTToken, error) {
	us, err := d.loadUserDetailByLogin(ctx, vm.Principal)
	if err != nil {
		return nil, errors.UserNotFound.Err(err).Descf("User not found by: %s", vm.Principal)
	}

	if !encrypt.ValidatePassword(us.Credentials, vm.Credentials) {
		return nil, errors.Unauthorized
	}

	if !us.Activated {
		return nil, errors.AccountNotActivated
	}

	return jwt.GenerateToken(jwt.Claims{
		Principal:   us.Principal,
		Authorities: us.Authorities,
	})
}

type LoginVM struct {
	Type        int8   `json:"type,omitempty"`
	Principal   string `json:"principal,omitempty" binding:"required"`
	Credentials string `json:"credentials,omitempty"`
}
