package karpenterreceiver

import (
    "context"
    "time"
    
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/receiver/receiverhelper"
)

func NewFactory() component.ReceiverFactory {
    return receiverhelper.NewFactory(
        "karpenter",
        createDefaultConfig,
        receiverhelper.WithMetrics(createMetricsReceiver),
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        ReceiverSettings: configmodels.ReceiverSettings{
            NameVal: "karpenter",
        },
        Endpoint: "http://localhost:8080/metrics",
        ScrapeInterval: 15 * time.Second,
    }
}
