// Code generated by Vanya: DO NOT EDIT.
// versions:
// 	vanya v0.0.0
// source: github.com/ivanmashin/my-service/configs/config.go

//go:generate go run github.com/ivanmashin/vanya/cmd/config
//go:build !vanya
// +build !vanya

package multiple_objs

import "github.com/ivanmashin/vanya/pkg/configs"

type Config struct {
	configs.Embedding

	HttpServerConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"http_server_config"`

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

func NewConfig(opts ...configs.Option) (Config, error) {
	c := NewDefaultConfig()

	err := c.Init(&c, opts...)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func NewDefaultConfig() Config {
	return Config{
		Embedding: configs.Embedding{},
		HttpServerConfig: struct {
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
			Database: "",
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
		OIDCConfig: struct {
			PartnerName      string `mapstructure:"partner_name"`
			ClientID         string `mapstructure:"client_id"`
			ClientSecret     string `mapstructure:"client_secret"`
			RedirectEndpoint string `mapstructure:"redirect_endpoint"`
		}{},
	}
}