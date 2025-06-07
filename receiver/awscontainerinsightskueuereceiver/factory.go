// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package awscontainerinsightskueuereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightskueuereceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	ci "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/containerinsight"
)

const (
	kueueMetricsStability = component.StabilityLevelDevelopment
)

var receiverType component.Type = component.MustNewType("awscontainerinsightskueuereceiver")

// Factory for awscontainerinsightskueuereceiver
const (
	// Rely on EC2 tags to auto-detect cluster name by default
	defaultClusterName = ""
)

// NewFactory creates a factory for AWS container insight receiver
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		receiverType,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, kueueMetricsStability))
}

// createDefaultConfig returns a default config for the receiver.
func createDefaultConfig() component.Config {
	return &Config{
		CollectionInterval: ci.DefaultCollectionInterval,
		ClusterName:        defaultClusterName,
	}
}

// CreateMetricsReceiver creates an AWS Container Insight receiver.
func createMetricsReceiver(
	_ context.Context,
	params receiver.Settings,
	baseCfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	rCfg := baseCfg.(*Config)
	return newAWSContainerInsightsKueueReceiver(params.TelemetrySettings, rCfg, consumer)
}
