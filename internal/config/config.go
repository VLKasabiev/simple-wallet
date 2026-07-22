package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/VLKasabiev/simple-wallet/pkg/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	IsProd   bool
	Web      *webParams
	Postgres *postgres.ConnectionData
	JWT      *JWTConfig 
}

type webParams struct {
	Port uint16
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn time.Duration
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

    viper.AutomaticEnv()

    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return &Config{
		IsProd: viper.GetBool("server.isProd"),
		Web: &webParams{
			Port: viper.GetUint16("server.port"),
		},
		Postgres: &postgres.ConnectionData{
			User:     viper.GetString("server.pg.user"),
			Password: viper.GetString("server.pg.password"),
			Host:     viper.GetString("server.pg.host"),
			Port:     viper.GetUint16("server.pg.port"),
			DBName:   viper.GetString("server.pg.database"),
			SSLMode:  viper.GetString("server.pg.sslmode"),
		},
		JWT: &JWTConfig{
			SecretKey: viper.GetString("server.jwt.secret"),
			ExpiresIn: viper.GetDuration("server.jwt.expires_in"),
		},
	}, nil
}

func (cfg *Config) GetWebPort() string {
	if cfg == nil || cfg.Web == nil {
		return ""
	}
	return strconv.Itoa(int(cfg.Web.Port))
}