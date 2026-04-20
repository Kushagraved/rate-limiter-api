package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Settings Settings `yaml:"settings"`
}

type Settings struct {
	Env         string         `yaml:"env"`
	Server      ServerCfg      `yaml:"server"`
	Logger      LoggerCfg      `yaml:"logger"`
	Redis       RedisCfg       `yaml:"redis"`
	RateLimiter RateLimiterCfg `yaml:"rate_limiter"`
}

type ServerCfg struct {
	Port string `yaml:"port"`
}

type LoggerCfg struct {
	Level string `yaml:"level"`
}

type RedisCfg struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// RateLimiterCfg controls rate limiting behaviour.
type RateLimiterCfg struct {
	MaxRequests int `yaml:"max_requests"`
	WindowSecs  int `yaml:"window_secs"`
}

func Load(path string) *Config {
	once.Do(func() {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read config %s: %v", path, err)
		}
		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("failed to parse config: %v", err)
		}
		instance = &cfg
	})
	return instance
}

func Get() *Config {
	if instance == nil {
		log.Fatal("config not loaded; call config.Load() first")
	}
	return instance
}
