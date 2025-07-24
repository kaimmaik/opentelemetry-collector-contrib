package kedareceiver

import (
    "time"
    "go.opentelemetry.io/collector/component"
)

// Config defines configuration for KEDA receiver
type Config struct {
    Endpoint       string        `mapstructure:"endpoint"`
    ScrapeInterval time.Duration `mapstructure:"scrape_interval"`
}

func (cfg *Config) Validate() error {
    return nil
}
