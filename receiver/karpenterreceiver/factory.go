package karpenterreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	promconfig "github.com/prometheus/prometheus/config"
	"github.com/prometheus/common/model"
)

const (
	typeStr         = "karpenterreceiver"
	defaultEndpoint = "http://karpenter.karpenter.svc.cluster.local:8000/metrics"
)

type Config struct {
	prometheusreceiver.Config `mapstructure:",squash"`
}

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelBeta),
	)
}

func createDefaultConfig() component.Config {
	promFactory := prometheusreceiver.NewFactory()
	promCfg := promFactory.CreateDefaultConfig().(*prometheusreceiver.Config)
	
	// Set up default Karpenter scrape config
	promCfg.PrometheusConfig.ScrapeConfigs = []*promconfig.ScrapeConfig{
		{
			JobName:        "karpenter-metrics",
			ScrapeInterval: model.Duration(30 * time.Second),
			ScrapeTimeout:  model.Duration(10 * time.Second),
			MetricsPath:    "/metrics",
			Scheme:         "http",
			StaticConfigs: []*promconfig.Group{
				{
					Targets: []model.LabelSet{
						{model.AddressLabel: model.LabelValue(defaultEndpoint)},
					},
				},
			},
		},
	}
	
	return &Config{Config: *promCfg}
}

func createMetricsReceiver(
	ctx context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	karpenterCfg := cfg.(*Config)
	promFactory := prometheusreceiver.NewFactory()
	return promFactory.CreateMetricsReceiver(ctx, set, &karpenterCfg.Config, nextConsumer)
}