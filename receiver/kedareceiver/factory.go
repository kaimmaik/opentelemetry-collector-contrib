package kedareceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const typeStr = "keda"

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelBeta),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Endpoint:       "keda-operator.keda-system.svc.cluster.local:8080",
		ScrapeInterval: 30 * time.Second,
	}
}

func createMetricsReceiver(
	ctx context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	return newKedaReceiver(cfg.(*Config), nextConsumer, set.Logger), nil
}