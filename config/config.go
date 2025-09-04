package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	AppName           = "interview-go-service"
	Host              = "localhost"
	Port              = "8000"
	CacheTTL          = time.Duration(time.Minute * 2)
	CacheClearTicker  = time.Duration(time.Second * 60)
	ApiRateLimitRate  = time.Duration(time.Second * 60)
	ApiRateLimitBurst = 10
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
		TTL         time.Duration `yaml:"ttl"`
		ClearTicker time.Duration `yaml:"clearticker"`
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
	if cfg.Cache.TTL == 0 {
		cfg.Cache.TTL = CacheTTL
	}
	if cfg.Cache.ClearTicker == 0 {
		cfg.Cache.ClearTicker = CacheClearTicker
	}
	if cfg.ApiRateLimit.Rate == 0 {
		cfg.ApiRateLimit.Rate = ApiRateLimitRate
	}
	if cfg.ApiRateLimit.Burst == 0 {
		cfg.ApiRateLimit.Burst = ApiRateLimitBurst
	}
}
