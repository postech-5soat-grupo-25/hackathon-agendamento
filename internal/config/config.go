package config

import (
	"github.com/spf13/viper"
)

var (
	Broker   *broker
	Postgres *postgres
)

func init() {
	viper.AutomaticEnv()
}

func LoadConfig() {
	LoadBrokerConfig()
	LoadPostgresConfig()
}
