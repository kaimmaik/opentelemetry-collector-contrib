module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kedareceiver

go 1.23

require (
	go.opentelemetry.io/collector v0.124.0
	go.opentelemetry.io/collector/component v0.124.0
	go.opentelemetry.io/collector/consumer v0.124.0
	go.opentelemetry.io/collector/receiver v0.124.0
	go.opentelemetry.io/collector/pdata v1.30.0
	go.uber.org/zap v1.27.0
	github.com/prometheus/client_golang v1.19.1
	github.com/prometheus/common v0.55.0
)