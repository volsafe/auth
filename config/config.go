package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DbName   string `json:"dbname"`
	}
	Jwt struct {
		Secret string `json:"secret"`
	}
	Server struct {
		Port string `json:"port"`
	}
}

var configInstance *Config

func LoadConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&configInstance); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}

func GetConfig() *Config {
	return configInstance
}
