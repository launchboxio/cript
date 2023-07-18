package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Store    string
	Redis    Redis
	Postgres Postgres
	Mysql    Mysql
	S3       S3
}

type Redis struct {
	Address  string
	Password string
	Database int
}

type Gorm struct {
	Url string
}

type Postgres struct {
	Gorm
}

type Mysql struct {
	Gorm
}

type S3 struct {
	Bucket string
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Store: "redis",
		Redis: Redis{
			Address: "redis-master.default:6379",
		},
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")

	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("postgres.url", "POSTGRES_URL")
	viper.BindEnv("mysql.url", "MYSQL_URL")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return cfg, nil
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	for _, key := range viper.AllKeys() {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		err := viper.BindEnv(key, envKey)
		if err != nil {
			return nil, err
		}
	}

	err := viper.Unmarshal(cfg)
	return cfg, err
}
