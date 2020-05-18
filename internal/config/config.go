package config

import (
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTP  HTTPConfig  `yaml:"http"`
	Redis RedisConfig `yaml:"redis"`
	Cache CacheConfig `yaml:"cache"`
	Service ServiceConfig `yaml:"service"`
}

type HTTPConfig struct {
  Host string `yaml:"host"`
	Port int `yaml:"port"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type CacheConfig struct {
	Size int `yaml:"size"`
}

type ServiceConfig struct {
  BaseURL string `yaml:"baseURL"`
}

func LoadConfig(filePath string) (*Config, error) {
	var cfg Config

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
  if err = decoder.Decode(&cfg); err != nil {
    return nil, err
  }

	return &cfg, nil
}
