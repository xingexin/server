package config

import "github.com/spf13/viper"

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
}

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
