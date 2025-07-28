package kedareceiver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, cfg.Validate())
}

func TestCreateMetricsReceiver(t *testing.T) {
	cfg := createDefaultConfig()
	metricsReceiver, err := createMetricsReceiver(
		context.Background(),
		receivertest.NewNopSettings(),
		cfg,
		consumertest.NewNop(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, metricsReceiver)
}

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	assert.NotNil(t, factory)
	assert.Equal(t, typeStr, factory.Type().String())
}