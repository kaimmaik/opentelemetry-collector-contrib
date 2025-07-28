package kedareceiver

import (
	"errors"
	"time"
)

type Config struct {
	Endpoint        string        `mapstructure:"endpoint"`
	ScrapeInterval  time.Duration `mapstructure:"scrape_interval"`
}

func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" {
		return errors.New("endpoint cannot be empty")
	}
	if cfg.ScrapeInterval <= 0 {
		return errors.New("scrape_interval must be positive")
	}
	return nil
}