package kedareceiver

import (
    "time"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/config/configmodels"
)

// Config defines configuration for KEDA receiver
type Config struct {
    configmodels.ReceiverSettings `mapstructure:",squash"`
    Endpoint       string        `mapstructure:"endpoint"`
    ScrapeInterval time.Duration `mapstructure:"scrape_interval"`
}
