package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	AppName = "interview-go-service"
	Host    = "localhost"
	Port    = "8000"
)

type Configuration struct {
	App struct {
		Name string `yaml:"name"`
	}

	Server struct {
		URL  string
		Port string
	}

	Cache struct {
		TTL         time.Duration `yaml:"ttl"`         // min
		ClearTicker time.Duration `yaml:"clearticker"` // sec
	} `yaml:"cache"`

	ApiRateLimit struct {
		Rate  time.Duration `yaml:"rate"`
		Burst int           `yaml:"burst"`
	} `yaml:"apiratelimit"`
}

func NewConfiguration(path string) (*Configuration, error) {
	cfg := &Configuration{}

	viper.SetConfigFile(path)
	_ = viper.ReadInConfig()
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}

	cfg.setDefaults()

	return cfg, nil
}

func (cfg *Configuration) setDefaults() {
	if cfg.App.Name == "" {
		cfg.App.Name = AppName
	}
	if cfg.Server.URL == "" {
		cfg.Server.URL = fmt.Sprintf("http://%s:%s", Host, Port)
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = Port
	}
}
