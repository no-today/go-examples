package setting

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// 不能直接通过配置lib读取配置, 必须通过该类读取配置, 屏蔽加载配置的细节, 需要更换配置lib的时候只需要改这里

var (
	Server  *serverConfig
	Async   *asyncConfig
	Mongodb *mongodbConfig
	Jwt     *jwtConfig
	Email   *emailConfig
	Redis   *redisConfig
)

const (
	// Server
	configServerPort         = "server.port"
	configServerEnv          = "server.env"
	configServerReadTimeout  = "server.read-timeout"
	configServerWriteTimeout = "server.write-timeout"
	configServerDomain       = "server.domain"

	// Async
	configAsyncWorkerPoolSize = "async.worker-pool-size"

	// Jwt
	configSecurityJwtBase64Secret           = "security.jwt.base64-secret"
	configSecurityJwtTokenValidityInSeconds = "security.jwt.token-validity-in-seconds"

	// Mongodb
	configDataMongoUri   = "data.mongo.uri"
	configDataMongoDebug = "data.mongo.debug"

	// Email
	configEmailHost     = "email.host"
	configEmailPort     = "email.port"
	configEmailUsername = "email.username"
	configEmailPassword = "email.password"

	// Redis
	configDataRedisAddr     = "data.redis.addr"
	configDataRedisUsername = "data.redis.username"
	configDataRedisPassword = "data.redis.password"
	configDataRedisDb       = "data.redis.db"
)

type serverConfig struct {
	Port         int32
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Domain       string
}

type asyncConfig struct {
	WorkerPoolSize int
}

type jwtConfig struct {
	Base64Secret           string
	TokenValidityInSeconds int64
}

type mongodbConfig struct {
	Uri   string
	Debug bool
}

type emailConfig struct {
	Host     string
	Port     int32
	Username string
	Password string
}

type redisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func Setup() {
	in := getAbsolutePath(2) + "/conf"
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(in)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Initial configuration failed, err: %v", err)
	}

	LoadConfigs()
	log.Println("Setting loading finished")
}

func getAbsolutePath(pre int) string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}

	for i := 0; i < pre; i++ {
		abPath = abPath[0:strings.LastIndex(abPath, string(os.PathSeparator))]
	}
	return abPath
}

func LoadConfigs() {
	loadServerConfig()
	loadAsyncConfig()
	loadMongodbConfig()
	loadJwtConfig()
	loadEmailConfig()
	loadRedisConfig()
}

func loadServerConfig() {
	Server = &serverConfig{
		Port:         viper.GetInt32(configServerPort),
		Env:          viper.GetString(configServerEnv),
		ReadTimeout:  time.Duration(viper.GetInt64(configServerReadTimeout)),
		WriteTimeout: time.Duration(viper.GetInt64(configServerWriteTimeout)),
		Domain:       viper.GetString(configServerDomain),
	}
}

func loadAsyncConfig() {
	Async = &asyncConfig{
		WorkerPoolSize: viper.GetInt(configAsyncWorkerPoolSize),
	}
}

func loadJwtConfig() {
	Jwt = &jwtConfig{
		Base64Secret:           viper.GetString(configSecurityJwtBase64Secret),
		TokenValidityInSeconds: viper.GetInt64(configSecurityJwtTokenValidityInSeconds),
	}
}

func loadMongodbConfig() {
	Mongodb = &mongodbConfig{
		Uri:   viper.GetString(configDataMongoUri),
		Debug: viper.GetBool(configDataMongoDebug),
	}
}

func loadEmailConfig() {
	Email = &emailConfig{
		Host:     viper.GetString(configEmailHost),
		Port:     viper.GetInt32(configEmailPort),
		Username: viper.GetString(configEmailUsername),
		Password: viper.GetString(configEmailPassword),
	}
}

func loadRedisConfig() {
	Redis = &redisConfig{
		Addr:     viper.GetString(configDataRedisAddr),
		Username: viper.GetString(configDataRedisUsername),
		Password: viper.GetString(configDataRedisPassword),
		DB:       viper.GetInt(configDataRedisDb),
	}
}
