package config

import (
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	Database   DatabaseConfig
	Serverport string
	JwtSecret  []byte
	CorsUrls   string
}

var CFG *Config

func LoadConfig() {
	viper.SetConfigFile("../.env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	CFG = &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		JwtSecret:  []byte("jwt_secret"),
		CorsUrls:   viper.GetString("CORS_URLS"),
		Serverport: viper.GetString("SERVER_PORT"),
	}
}
