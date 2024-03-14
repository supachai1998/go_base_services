package configs

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// Config contains all server configurations.
type Config struct {
	ENV       string
	Version   string
	DebugMode bool
	Server    struct {
		Host            string
		Port            int
		ShutdownTimeout time.Duration
		RequestTimeout  time.Duration
	}

	DB    DBConfig
	Redis Redis

	PrettyLog bool

	BaseUrl string

	// swagger info
	SwaggerTitle       string
	SwaggerDescription string
	SwaggerContact     SwaggerContact
	SwaggerLicense     SwaggerLicense

	// Echo data
	ECHODATA struct {
		REQ bool
		RES bool
	}

	// Gorm debug
	DbDebug bool

	// Auth
	VerifyTokenExpire time.Duration

	AdminAuth AuthConfig
	UserAuth  AuthConfig

	Passphrase string

	// Cache Expire
	CacheExpireStaff time.Duration
}

type SwaggerContact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type SwaggerLicense struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type DBConfig struct {
	URI        string
	HOST       string
	PORT       string
	USER       string
	PASSWORD   string
	DBNAME     string
	CHARSET    string
	SSLMODE    string
	TIMEZONE   string
	DBDEBUG    bool
	CACHEDEBUG bool
}

type Redis struct {
	HOST      string
	DB        int
	PASSWORD  string
	TLSSERVER string
}

type Cron struct {
	Enables []string
}
type AuthConfig struct {
	JWTSecret                 string
	AccessTokenDuration       time.Duration
	RefreshTokenDuration      time.Duration
	VerifyTokenDuration       time.Duration
	AccountLockoutMaxAttempts int
	ResendOTPMaxAttempts      int
	VerifyOTPMaxAttempts      int
	ReturnOTP                 bool
	DemoUser                  struct {
		Email string
		Tel   string
		Pin   string
	}
	LenTempPwd int
}

var (
	_, b, _, _ = runtime.Caller(0)
	BasePath   = filepath.Dir(b)
	Root       = filepath.Join(filepath.Dir(b), "../.")
	ModeDev    = "development"
	ModeProd   = "production"
	ModeTest   = "test"
)

func (cfg DBConfig) GetURI() string {
	uri := cfg.URI
	if uri == "" {
		uri = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			cfg.HOST, cfg.USER, cfg.PASSWORD, cfg.DBNAME, cfg.PORT, cfg.SSLMODE, cfg.TIMEZONE)
	}
	return uri
}

func (cfg Redis) GetOptions() *redis.Options {
	return &redis.Options{
		Addr:     cfg.HOST,
		Password: cfg.PASSWORD,
		DB:       cfg.DB,
	}
}

func (cfg Redis) GetURI() string {
	uri := cfg.HOST
	if uri == "" {
		uri = fmt.Sprintf("redis://%s", cfg.HOST)
	}
	return uri
}

func (cfg Redis) GetTLSServer() string {
	return cfg.TLSSERVER
}

func ParseConfig(mode ...string) (cfg *Config, err error) {
	viper.AddConfigPath(Root + "/configs")
	MODE := ""
	if len(mode) > 0 {
		MODE = fmt.Sprintf(".%s", mode[0])
	}

	// configs
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// secret
	viper.SetConfigName("secret" + MODE)
	viper.SetConfigType("yaml")
	err = viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	// versioning
	viper.SetConfigName("version")
	viper.SetConfigType("yaml")
	err = viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	// swagger info
	viper.SetConfigName("swagger")
	viper.SetConfigType("yaml")
	err = viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
