package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBConnectionString string
	DBName             string
	JWTSecret          string
}

var Cfg Config

func LoadConfig() {
	viper.SetConfigFile(".env") // load from .env
	viper.AddConfigPath(".")    // path to look for config file

	// Optional: read from environment variables too
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err := viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	fmt.Println("✅ Config loaded successfully")
}
