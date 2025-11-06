package config

import (
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server struct {
		Port int
	}
	DataBase struct {
		Driver string
		DSN    string
	}
	Logger struct {
		Level string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
		PoolSize int
	}
}

// LoadConfig 从 config.yaml 加载配置文件
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
