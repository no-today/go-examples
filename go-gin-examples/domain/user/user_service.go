package user

import (
	"cathub.me/go-web-examples/middleware/security"
	"cathub.me/go-web-examples/pkg/async"
	"cathub.me/go-web-examples/pkg/data"
	"cathub.me/go-web-examples/pkg/database"
	"cathub.me/go-web-examples/pkg/errors"
	"cathub.me/go-web-examples/pkg/setting"
	"cathub.me/go-web-examples/pkg/utils/encrypt"
	"cathub.me/go-web-examples/pkg/utils/mail"
	"cathub.me/go-web-examples/pkg/utils/strutil"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"sync"
	"time"
)

// _userService Singleton
var _userService UserService
var _onceUserService sync.Once

func GetUserService() UserService {
	_onceUserService.Do(func() {
		_userService = &userService{userRepository: GetUserRepository(), redis: database.GetRedisClient()}
	})
	return _userService
}

const (
	// Redis key, user activation code
	rkUserActivationCode = "user_activation_code"
)

type UserService interface {
	Register(ctx context.Context, dto RegisterUserDTO) (*UserDTO, error)
	ResendActivateEmail(ctx context.Context, email string) error
	Activation(ctx context.Context, code string) error
	ClearUnActivationUser(ctx context.Context) (count int64)
	FindAll(ctx context.Context, pageable data.Pageable) (*data.Pageable, []*UserDTO, error)
}

type userService struct {
	userRepository UserRepository
	redis          *redis.Client
}

func (u *userService) Register(ctx context.Context, dto RegisterUserDTO) (*UserDTO, error) {
	us, err := u.userRepository.FindByUsername(ctx, dto.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if us != nil {
		return nil, errors.UsernameAlreadyExists.Descf("Username: %s, already exists.", dto.Username)
	}

	us, err = u.userRepository.FindByEmail(ctx, dto.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if us != nil {
		return nil, errors.EmailAlreadyExists.Descf("Email: %s, already exists.", dto.Email)
	}

	passwordHash, _ := encrypt.Encrypt(dto.Password)

	newUser := &User{
		Id:               primitive.NewObjectID(),
		Username:         dto.Username,
		Password:         passwordHash,
		Roles:            []string{security.USER},
		Email:            dto.Email,
		Activated:        false,
		CreatedDate:      time.Now(),
		LastModifiedDate: time.Now(),
	}

	if err = u.userRepository.Insert(ctx, newUser); err != nil {
		return nil, err
	}
	// Shouldn't return password to client
	newUser.Password = ""

	// Send activate email, cannot log in before activation
	u.sendActivateEmail(ctx, *newUser)

	return ToDTO(newUser), nil
}

func (u *userService) ResendActivateEmail(ctx context.Context, email string) error {
	us, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.EmailNotFound.Descf("Email not found by: %s", email)
		}
		return err
	}

	if us.Activated {
		return errors.UserAlreadyActivated.Descf("User already activated, please log in")
	}

	u.sendActivateEmail(ctx, *us)
	return nil
}

func (u *userService) sendActivateEmail(ctx context.Context, user User) {
	code := strings.ToUpper(strutil.RandStringRunes(32))
	u.setUserActivationCode(user.Id.Hex(), code)

	// 异步执行
	async.Submit(func() {
		if err := mail.SendEmail("welcome register go-web-examples",
			fmt.Sprintf("activate on click: %s/api/v1/activation/%s", setting.Server.Domain, code),
			user.Email); err != nil {
			log.Err(ctx.Err()).Str("email", user.Email).Str("code", code).Msgf("Send activate email failed")
		} else {
			log.Info().Str("email", user.Email).Str("code", code).Msgf("Send activate email success")
		}
	})
}

func getUserActivationCodeKey(code string) string {
	return fmt.Sprintf("%s:%s", rkUserActivationCode, code)
}

func (u *userService) Activation(ctx context.Context, code string) error {
	result, err := u.getUserActivationCode(code)
	if err != nil {
		if err == redis.Nil {
			return errors.InvalidActivateCode
		}
		return err
	}

	id, _ := primitive.ObjectIDFromHex(result)

	if _, err = u.userRepository.UpdateById(ctx, &User{
		Id:               id,
		Activated:        true,
		ActivatedDate:    time.Now(),
		LastModifiedDate: time.Now(),
	}); err != nil {
		return err
	}

	u.deleteUserActivationCode(code)

	return nil
}

func (u *userService) setUserActivationCode(id, code string) *redis.StatusCmd {
	return u.redis.SetEX(context.Background(), getUserActivationCodeKey(code), id, 72*time.Hour)
}

func (u *userService) getUserActivationCode(code string) (string, error) {
	result, err := u.redis.Get(context.Background(), getUserActivationCodeKey(code)).Result()
	return result, err
}

func (u *userService) deleteUserActivationCode(code string) *redis.IntCmd {
	return u.redis.Del(context.Background(), fmt.Sprintf("user_activation_code:%s", code))
}

func (u *userService) ClearUnActivationUser(ctx context.Context) (count int64) {
	// Default clear 3 days not activation user
	return u.clearUnActivationUser(ctx, time.Now().Add(-72*time.Hour))
}

func (u *userService) clearUnActivationUser(ctx context.Context, lteTime time.Time) (count int64) {
	page := int64(1)
	size := int64(1000)

	for true {
		users, err := u.userRepository.FindAllByActivatedIsFalseAndCreatedDateLessThan(ctx, lteTime, data.Pageable{Page: page, Size: size, Fields: []string{"_id", "username"}})
		if err != nil {
			log.Err(err).Msg("Clear un activation user failed")
			return count
		}

		if users == nil {
			return count
		}

		for _, user := range users {
			_ = u.userRepository.DeleteById(ctx, user.Id)
			count++
		}

		page++
	}

	return count
}

func (u *userService) FindAll(ctx context.Context, pageable data.Pageable) (*data.Pageable, []*UserDTO, error) {
	pageInfo, users, err := u.userRepository.FindAll(ctx, pageable)
	return pageInfo, ToDTOS(users), err
}

// RegisterUserDTO 注册用户业务模型
type RegisterUserDTO struct {
	Username string `json:"username,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
	Email    string `json:"email,omitempty" binding:"required"`
}
