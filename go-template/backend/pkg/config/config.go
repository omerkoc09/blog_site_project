package config

import (
	"time"

	"github.com/spf13/viper"
)

/*
TODO: yeni bir config eklerken setDefaults ve bindEnvs e eklemeyi unutma !
*/
var config *Config

type Config struct {
	IsDevelopment bool
	Server        ServerConfig
	Database      DbConfig
	Redis         RedisConfig
}

type ServerConfig struct {
	Port         int
	JwtSecret    string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	LogPath      string
	LogLevel     string
	UploadDir    string

	JwtAccessTokenExpireMinute time.Duration
	JwtRefreshTokenExpireHour  time.Duration
}

type DbConfig struct {
	Name        string
	Username    string
	Password    string
	Host        string
	Port        string
	Debug       bool
	MaxPoolSize int
	MaxIdleConn int
	MaxLifetime int
}

type RedisConfig struct {
	Host         string
	Password     string
	Db           int
	WriteTimeout time.Duration
}

func setDefaults() {
	viper.SetDefault("isDevelopment", false)
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("server.logPath", "./log")
	viper.SetDefault("server.logLevel", "debug")
	viper.SetDefault("server.jwtSecret", "gizli-bir-jwt-secreti-olacak-ins")
	viper.SetDefault("server.jwtAccessTokenExpireMinute", 5)
	viper.SetDefault("server.jwtRefreshTokenExpireHour", 240)
	viper.SetDefault("server.readTimeout", 30)
	viper.SetDefault("server.writeTimeout", 30)
	viper.SetDefault("server.IdleTimeout", 120)
	viper.SetDefault("server.uploadDir", "./upload")

	viper.SetDefault("database.name", "go_db")
	viper.SetDefault("database.username", "omerkoc")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.debug", true)
	viper.SetDefault("database.maxPoolSize", 5)
	viper.SetDefault("database.maxIdleConn", 1)
	viper.SetDefault("database.maxLifetime", 1800)

	viper.SetDefault("redis.host", "data.hayteknoloji.com:6379")
	viper.SetDefault("redis.password", "JNfzkKwXdfmtnbdZsziQxMUyhVJF7uB6LMVwCPwfbjSXOSVlfi1xa5EE6LpD8NkP7Ud4J0zdjovHW2Yt")
	viper.SetDefault("redis.db", 9)
	viper.SetDefault("redis.writeTimeout", time.Second*5)
}

// os environment'ından okumak için. şimdilik çokda lazım değil..
func bindEnvs() {
	bindEnv := func(input ...string) {
		err := viper.BindEnv(input...)
		if err != nil {
			panic(err)
		}
	}

	bindEnv("isDevelopment", "IS_DEVELOPMENT")
	bindEnv("server.port", "SERVER_PORT")
	bindEnv("server.logPath", "SERVER_LOG_PATH")
	bindEnv("server.logLevel", "SERVER_LOG_LEVEL")
	bindEnv("server.jwtSecret", "SERVER_JWT_SECRET")
	bindEnv("server.jwtAccessTokenExpireMinute", "SERVER_JWT_ACCESS_TOKEN_EXPIRE_MINUTE")
	bindEnv("server.jwtRefreshTokenExpireHour", "SERVER_JWT_REFRESH_TOKEN_EXPIRE_HOUR")
	bindEnv("server.timeout", "SERVER_TIMEOUT")
	bindEnv("server.readTimeout", "SERVER_READ_TIMEOUT")
	bindEnv("server.writeTimeout", "SERVER_WRITE_TIMEOUT")
	bindEnv("server.IdleTimeout", "SERVER_IDLE_TIMEOUT")
	bindEnv("server.uploadDir", "SERVER_UPLOAD_DIR")

	bindEnv("database.name", "DB_NAME")
	bindEnv("database.username", "DB_USERNAME")
	bindEnv("database.password", "DB_PASSWORD")
	bindEnv("database.host", "DB_HOST")
	bindEnv("database.port", "DB_PORT")
	bindEnv("database.debug", "DB_DEBUG")
	bindEnv("database.maxPoolSize", "DB_MAX_POOL_SIZE")
	bindEnv("database.maxIdleConn", "DB_MAX_IDLE_CONN")
	bindEnv("database.maxLifetime", "DB_LIFETIME")

	bindEnv("redis.host", "REDIS_HOST")
	bindEnv("redis.password", "REDIS_PASSWORD")
	bindEnv("redis.db", "REDIS_DB")
	bindEnv("redis.writeTimeout", "REDIS_WRITE_TIMEOUT")

}

func Setup() (*Config, error) {
	setDefaults()
	bindEnvs()
	// Auto read env variables
	viper.AutomaticEnv()

	// Unmarshal config file to struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

func Get() *Config {
	if config == nil {
		panic("Config is not initialized")
	}
	return config
}
