package kedareceiver

import (
    "context"
    "time"
    
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/receiver/receiverhelper"
)

func NewFactory() component.ReceiverFactory {
    return receiverhelper.NewFactory(
        "keda",
        createDefaultConfig,
        receiverhelper.WithMetrics(createMetricsReceiver),
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        ReceiverSettings: configmodels.ReceiverSettings{
            NameVal: "keda",
        },
        Endpoint: "keda-operator.keda:8080",
        ScrapeInterval: 15 * time.Second,
    }
}
