package config

import (
	"github.com/spf13/viper"
)

var (
	Broker *broker
)

func init() {
	viper.AutomaticEnv()
}

func LoadConfig() {
	LoadBrokerConfig()
}
