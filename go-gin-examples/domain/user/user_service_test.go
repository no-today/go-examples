package user

import (
	"cathub.me/go-gin-examples/middleware/security"
	"cathub.me/go-gin-examples/pkg/database"
	"cathub.me/go-gin-examples/pkg/errors"
	"cathub.me/go-gin-examples/pkg/logging"
	"cathub.me/go-gin-examples/pkg/setting"
	"cathub.me/go-gin-examples/pkg/utils/strutil"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func init() {
	logging.Setup()
	setting.Setup()
}

func TestGetUserService(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			"#",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserService(); got == nil {
				t.Errorf("GetUserService() = %v, want not nil", got)
			}
		})
	}
}

func Test_getUserActivationCodeKey(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "0000",
			args: args{code: "0000"},
			want: rkUserActivationCode + ":0000",
		},
		{
			name: "0001",
			args: args{"0001"},
			want: rkUserActivationCode + ":0001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUserActivationCodeKey(tt.args.code); got != tt.want {
				t.Errorf("getUserActivationCodeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Activation(t *testing.T) {
	type fields struct {
		userRepository UserRepository
		redis          *redis.Client
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		before  func(service *userService)
		after   func(service *userService)
	}{
		// TODO: Add test cases.
		{
			name: "InvalidActivateCode",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			before:  nil,
			args:    args{code: "0000"},
			wantErr: true,
		},
		{
			name: "ActivatedSuccessfully",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			args:    args{code: "code-01"},
			wantErr: false,
			before: func(service *userService) {
				service.setUserActivationCode(primitive.NewObjectID().Hex(), "code-01")
			},
			after: func(service *userService) {
				service.deleteUserActivationCode("code-01")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				userRepository: tt.fields.userRepository,
				redis:          tt.fields.redis,
			}

			if tt.before != nil {
				tt.before(u)
			}

			if err := u.Activation(context.Background(), tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("Activation() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.after != nil {
				tt.after(u)
			}
		})
	}
}

func Test_userService_clearUnActivationUser(t *testing.T) {
	type fields struct {
		userRepository UserRepository
		redis          *redis.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
		before func(repository UserRepository)
		after  func(repository UserRepository)
	}{
		// TODO: Add test cases.
		{
			name: "NotMatcher",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			want: 0,
		},
		{
			name: "Matched",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			want: 100,
			before: func(repository UserRepository) {
				var users []*User
				for i := 0; i < 100; i++ {
					users = append(users, &User{
						Id:               primitive.NewObjectID(),
						Username:         fmt.Sprintf("%s-%d", strutil.RandStringRunes(10), i),
						Password:         "",
						Roles:            []string{security.USER},
						Email:            fmt.Sprintf("%s-%d@xx.com", strutil.RandStringRunes(10), i),
						Activated:        false,
						CreatedDate:      time.Now(),
						LastModifiedDate: time.Now(),
					})
				}
				_ = repository.InsertAll(context.Background(), users)
			},
			after: func(repository UserRepository) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				userRepository: tt.fields.userRepository,
				redis:          tt.fields.redis,
			}

			if tt.before != nil {
				tt.before(u.userRepository)
			}

			if got := u.clearUnActivationUser(context.Background(), time.Now()); got != tt.want {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}

			if tt.after != nil {
				tt.after(u.userRepository)
			}
		})
	}
}

func Test_userService_Register(t *testing.T) {
	type fields struct {
		userRepository UserRepository
		redis          *redis.Client
	}
	type args struct {
		dto RegisterUserDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
		before  func(repository UserRepository)
		after   func(repository UserRepository)
	}{
		// TODO: Add test cases.
		{
			name: "Success",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			args: args{dto: RegisterUserDTO{
				Username: "register-user-01",
				Password: "xxx",
				Email:    "chengbdtb@163.com",
			}},
			want:    "register-user-01",
			wantErr: nil,
			after:   deleteByUsername("register-user-01"),
		},
		{
			name: "UsernameAlreadyExists",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			args: args{dto: RegisterUserDTO{
				Username: "register-user-02",
				Password: "xxx",
				Email:    "xxxx@163.com",
			}},
			want:    "register-user-02",
			wantErr: errors.UsernameAlreadyExists,
			before: func(repository UserRepository) {
				_ = repository.Insert(context.Background(), &User{
					Id:               primitive.NewObjectID(),
					Username:         "register-user-02",
					Password:         "",
					Roles:            []string{security.USER},
					Email:            "xxxx@163.com",
					Activated:        false,
					CreatedDate:      time.Now(),
					LastModifiedDate: time.Now(),
				})
			},
			after: deleteByUsername("register-user-02"),
		},
		{
			name: "EmailAlreadyExists",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			args: args{dto: RegisterUserDTO{
				Username: "register-user-03",
				Password: "xxx",
				Email:    "chengbdtb@163.com",
			}},
			want:    "",
			wantErr: errors.EmailAlreadyExists,
			before: func(repository UserRepository) {
				_ = repository.Insert(context.Background(), &User{
					Id:               primitive.NewObjectID(),
					Username:         "register-user-xxx",
					Password:         "",
					Roles:            []string{security.USER},
					Email:            "chengbdtb@163.com",
					Activated:        false,
					CreatedDate:      time.Now(),
					LastModifiedDate: time.Now(),
				})
			},
			after: deleteByUsername("register-user-xxx"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				userRepository: tt.fields.userRepository,
				redis:          tt.fields.redis,
			}

			if tt.before != nil {
				tt.before(u.userRepository)
			}

			got, err := u.Register(context.Background(), tt.args.dto)

			if err != nil || tt.wantErr != nil {
				if tt.wantErr == nil {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					problem, ok := err.(errors.Problem)
					if !ok || problem.Title != tt.wantErr.(errors.Problem).Title {
						t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
					}
				}
			} else {
				if got.Username != tt.want {
					t.Errorf("Register() got = %v, want %v", got.Username, tt.want)
				}
			}

			if tt.after != nil {
				tt.after(u.userRepository)
			}
		})
	}
}

func deleteByUsername(username string) func(repository UserRepository) {
	return func(repository UserRepository) {
		us, _ := repository.FindByUsername(context.Background(), username)
		_ = repository.DeleteById(context.Background(), us.Id)
	}
}

func Test_userService_ResendActivateEmail(t *testing.T) {
	type fields struct {
		userRepository UserRepository
		redis          *redis.Client
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		before  func(repository UserRepository)
		after   func(repository UserRepository)
	}{
		// TODO: Add test cases.
		{
			name: "NotFound",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			args:    args{email: "no.today@outlook.com"},
			wantErr: errors.EmailNotFound,
		},
		{
			name: "AlreadyActivated",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			args: args{
				email: "already_activated@outlook.com",
			},
			wantErr: errors.UserAlreadyActivated,
			before: func(repository UserRepository) {
				_ = repository.Insert(context.Background(), &User{
					Id:               primitive.NewObjectID(),
					Username:         strutil.RandStringRunes(32),
					Password:         "",
					Roles:            []string{security.USER},
					Email:            "already_activated@outlook.com",
					Activated:        true,
					ActivatedDate:    time.Now(),
					CreatedDate:      time.Now(),
					LastModifiedDate: time.Now(),
				})
			},
			after: func(repository UserRepository) {
				us, _ := repository.FindByEmail(context.Background(), "already_activated@outlook.com")
				_ = repository.DeleteById(context.Background(), us.Id)
			},
		},
		{
			name: "Success",
			fields: fields{
				userRepository: GetUserRepository(),
				redis:          database.GetRedisClient(),
			},
			args: args{
				email: "chengbdtb@163.com",
			},
			wantErr: nil,
			before: func(repository UserRepository) {
				_ = repository.Insert(context.Background(), &User{
					Id:               primitive.NewObjectID(),
					Username:         strutil.RandStringRunes(32),
					Password:         "",
					Roles:            []string{security.USER},
					Email:            "chengbdtb@163.com",
					Activated:        false,
					CreatedDate:      time.Now(),
					LastModifiedDate: time.Now(),
				})
			},
			after: func(repository UserRepository) {
				us, _ := repository.FindByEmail(context.Background(), "chengbdtb@163.com")
				_ = repository.DeleteById(context.Background(), us.Id)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				userRepository: tt.fields.userRepository,
				redis:          tt.fields.redis,
			}

			if tt.before != nil {
				tt.before(u.userRepository)
			}

			if err := u.ResendActivateEmail(context.Background(), tt.args.email); err != nil || tt.wantErr != nil {
				if tt.wantErr == nil {
					t.Errorf("ResendActivateEmail() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					problem, ok := err.(errors.Problem)
					if !ok || problem.Title != tt.wantErr.(errors.Problem).Title {
						t.Errorf("ResendActivateEmail() error = %v, wantErr %v", err, tt.wantErr)
					}
				}
			}

			if tt.after != nil {
				tt.after(u.userRepository)
			}
		})
	}
}
