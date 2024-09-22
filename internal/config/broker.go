package config

import "github.com/spf13/viper"

const (
	envHost = "BROKER_HOST"

	defaultHost = "localhost:9092"
)

type broker struct {
	host string
}

func LoadBrokerConfig() {
	viper.SetDefault(envHost, defaultHost)
	Broker = &broker{
		host: viper.GetString(envHost),
	}
}
