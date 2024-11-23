package config

import "github.com/ShatteredRealms/go-common-service/pkg/config"

var (
	Version = "v1.0.0"
)

type InventoryConfig struct {
	config.BaseConfig `yaml:",inline" mapstructure:",squash"`
	Postgres          config.DBPoolConfig `yaml:"postgres"`
}

func NewInventoryConfig() *InventoryConfig {
	return &InventoryConfig{
		BaseConfig: config.BaseConfig{
			Server: config.ServerAddress{
				Host: "localhost",
				Port: "8084",
			},
			Keycloak: config.KeycloakConfig{
				BaseURL:      "localhost:8080",
				Realm:        "default",
				Id:           "5114c8a3-5035-4f44-ba55-5d0048ea40ef",
				ClientId:     "sro-inventory-service",
				ClientSecret: "**********",
			},
			Mode:                "local",
			LogLevel:            0,
			OpenTelemtryAddress: "localhost:4317",
		},
		Postgres: config.DBPoolConfig{
			Master: config.DBConfig{
				ServerAddress: config.ServerAddress{},
				Name:          "inventory-service",
				Username:      "postgres",
				Password:      "password",
			},
		},
	}
}
