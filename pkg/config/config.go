package config

import (
	"context"

	cconfig "github.com/ShatteredRealms/go-common-service/pkg/config"
)

var (
	Version = "v1.0.0"
)

type InventoryConfig struct {
	cconfig.BaseConfig `yaml:",inline" inventorystructure:",squash"`
	Postgres           cconfig.DBPoolConfig `yaml:"postgres"`
}

func NewInventoryConfig(ctx context.Context) (*InventoryConfig, error) {
	config := &InventoryConfig{
		BaseConfig: cconfig.BaseConfig{
			Server: cconfig.ServerAddress{
				Host: "localhost",
				Port: "8085",
			},
			Keycloak: cconfig.KeycloakConfig{
				BaseURL:      "localhost:8080",
				Realm:        "default",
				Id:           "7b575e9b-c687-4cdc-b210-67c59b5f380f",
				ClientId:     "sro-inventory-service",
				ClientSecret: "**********",
			},
			Mode:                "local",
			LogLevel:            0,
			OpenTelemtryAddress: "localhost:4317",
		},
		Postgres: cconfig.DBPoolConfig{
			Master: cconfig.DBConfig{
				ServerAddress: cconfig.ServerAddress{},
				Name:          "inventory-service",
				Username:      "postgres",
				Password:      "password",
			},
		},
	}

	err := cconfig.BindConfigEnvs(ctx, "sro-inventory", config)
	return config, err
}
