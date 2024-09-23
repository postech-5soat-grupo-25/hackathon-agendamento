package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	envDBHost   = "DB_HOST"
	envDBPort   = "DB_PORT"
	envDBUser   = "DB_USER"
	envDBPasswd = "DB_PASSWD"
	envDBName   = "DB_NAME"

	defaultDBHost   = "localhost"
	defaultDBPort   = "5432"
	defaultDBUser   = "postgres"
	defaultDBPasswd = "postgres"
	defaultDBName   = "agendamento"
)

type postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func GetPostgresConnString() string {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString(envDBHost),
		viper.GetString(envDBPort),
		viper.GetString(envDBUser),
		viper.GetString(envDBPasswd),
		viper.GetString(envDBName))
	return connString
}

func LoadPostgresConfig() {
	viper.SetDefault(envDBHost, defaultDBHost)
	viper.SetDefault(envDBPort, defaultDBPort)
	viper.SetDefault(envDBUser, defaultDBUser)
	viper.SetDefault(envDBPasswd, defaultDBPasswd)
	viper.SetDefault(envDBName, defaultDBName)

	Postgres = &postgres{
		Host:     viper.GetString(envDBHost),
		Port:     viper.GetString(envDBPort),
		User:     viper.GetString(envDBUser),
		Password: viper.GetString(envDBPasswd),
		Dbname:   viper.GetString(envDBName),
	}
}
