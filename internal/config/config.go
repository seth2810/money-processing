package config

import (
	"fmt"

	"github.com/iamolegga/enviper"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"db_host"`
	User     string `mapstructure:"db_user"`
	Password string `mapstructure:"db_password"`
	Name     string `mapstructure:"db_name"`
	Port     uint16 `mapstructure:"db_port"`
}

type HTTPConfig struct {
	Host, Port string
}

type Config struct {
	Database DatabaseConfig `mapstructure:",squash"`
	Server   HTTPConfig     `mapstructure:",squash"`
}

func Read() (*Config, error) {
	var cfg Config

	v := enviper.New(viper.New())

	v.AddConfigPath(".")
	v.SetConfigType("env")
	v.SetConfigName("config")

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("while unmarshal config: %w", err)
	}

	return &cfg, nil
}
