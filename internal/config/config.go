package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBName  string `mapstructure:"DATABASE_NAME"`
	DBPass  string `mapstructure:"DATABASE_PASSWORD"`
	DBUser  string `mapstructure:"DATABASE_USER"`
	DBHost  string `mapstructure:"DATABASE_HOST"`
	DBPort  string `mapstructure:"DATABASE_PORT"`
	AppPort string `mapstructure:"PORT"`
}

func MustLoad() *Config {
	viper.SetConfigFile("/src/.env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot read the config: " + err.Error())
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)

	if err != nil {
		panic("error decoding cfg: " + err.Error())
	}
	requiredFields := []string{"DATABASE_NAME", "DATABASE_PASSWORD", "DATABASE_USER", "DATABASE_HOST", "DATABASE_PORT", "PORT"}
	for _, field := range requiredFields {
		if !viper.IsSet(field) {
			panic("Required enviroment string doesnt exist: " + field)
		}
	}

	return &cfg
}
