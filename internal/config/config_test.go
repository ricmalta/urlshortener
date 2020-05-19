package config_test

import (
	"os"
	"testing"

	"github.com/ricmalta/urlshortner/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestInvalidFile(t *testing.T) {
	const invalidPath = "./config_invalid.yaml"
	cfg, err := config.LoadConfig(invalidPath)
	assert.NotNil(t, err, "returns an error")
	assert.Nil(t, cfg, "config instance should be null")
}

func TestValidFile(t *testing.T) {
	const validPath = "./config.yaml"
	cfg, err := config.LoadConfig(validPath)
	assert.Nil(t, err, "return no error")
	assert.NotNil(t, cfg, "returns a valid config instance")

	var testCfg config.Config
	file, err := os.Open(validPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&testCfg); err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, cfg.HTTP.Host, testCfg.HTTP.Host)
	assert.Equal(t, cfg.HTTP.Port, testCfg.HTTP.Port)
	assert.Equal(t, cfg.Cache.Size, testCfg.Cache.Size)
	assert.Equal(t, cfg.Redis.Host, testCfg.Redis.Host)
	assert.Equal(t, cfg.Redis.Database, testCfg.Redis.Database)
	assert.Equal(t, cfg.Redis.Password, testCfg.Redis.Password)
	assert.Equal(t, cfg.Service.BaseURL, testCfg.Service.BaseURL)
	assert.Equal(t, cfg.Cache.Size, testCfg.Cache.Size)
	assert.Equal(t, cfg.LogLevel, testCfg.LogLevel)
}
