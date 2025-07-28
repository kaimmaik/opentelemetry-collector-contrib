package karpenterreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const typeStr = "karpenter"

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelBeta),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Endpoint:       "karpenter.karpenter.svc.cluster.local:8000",
		ScrapeInterval: 30 * time.Second,
	}
}

func createMetricsReceiver(
	ctx context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	return newKarpenterReceiver(cfg.(*Config), nextConsumer, set.Logger), nil
}