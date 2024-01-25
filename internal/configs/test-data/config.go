//go:build vanya
// +build vanya

package test_data

import (
	base "github.com/ivanmashin/vanya"
	"github.com/ivanmashin/vanya/pkg/configs"
)

func main() {
	base.BuildConfigs(
		configs.HttpServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		configs.GrpcServerConfig{
			Host: "localhost",
			Port: "1000",
		},
		configs.PostgresConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "admin",
			Password: "admin",
			Database: "",
		},
		configs.RedisConfig{
			Host: "localhost",
			Port: "6379",
			DB:   1,
		},
		OIDCConfig{},
	)
}

type OIDCConfig struct {
	PartnerName      string
	ClientID         string
	ClientSecret     string
	RedirectEndpoint string
}

type Config struct {
	configs.Embedding

	HttpConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"http_config"`

	GrpcServerConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"grpc_server_config"`

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	} `mapstructure:"postgres_config"`

	RedisConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis_config"`

	OIDCConfig struct {
		PartnerName      string `mapstructure:"partner_name"`
		ClientID         string `mapstructure:"client_id"`
		ClientSecret     string `mapstructure:"client_secret"`
		RedirectEndpoint string `mapstructure:"redirect_endpoint"`
	} `mapstructure:"oidc_config"`
}

func NewDefaultConfig() Config {
	return Config{
		Embedding: configs.Embedding{},
		HttpConfig: struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		}{
			Host: "localhost",
			Port: "8080",
		},
		GrpcServerConfig: struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		}{
			Host: "localhost",
			Port: "1000",
		},
		PostgresConfig: struct {
			Host     string `mapstructure:"host"`
			Port     string `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Database string `mapstructure:"database"`
		}{
			Host:     "localhost",
			Port:     "5432",
			User:     "admin",
			Password: "admin",
		},
		RedisConfig: struct {
			Host     string `mapstructure:"host"`
			Port     string `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			DB       int    `mapstructure:"db"`
		}{
			Host: "localhost",
			Port: "6379",
			DB:   1,
		},
	}
}
