package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTP  HTTPConfig
	Redis RedisConfig
	Cache CacheConfig
}

type HTTPConfig struct {
	Port int
}

type RedisConfig struct {
	Host     string
	Password string
	Database int
}

type CacheConfig struct {
	Size int
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.AutomaticEnv()

	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
