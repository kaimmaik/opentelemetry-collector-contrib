package karpenterreceiver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type karpenterReceiver struct {
	cfg          *Config
	consumer     consumer.Metrics
	logger       *zap.Logger
	cancel       context.CancelFunc
}

func newKarpenterReceiver(set receiver.Settings, cfg *Config, nextConsumer consumer.Metrics) (receiver.Metrics, error) {
	return &karpenterReceiver{
		cfg:      cfg,
		consumer: nextConsumer,
		logger:   set.Logger,
	}, nil
}

func (r *karpenterReceiver) Start(ctx context.Context, host component.Host) error {
	ctx, r.cancel = context.WithCancel(ctx)
	
	go r.scrapeLoop(ctx)
	
	return nil
}

func (r *karpenterReceiver) Shutdown(ctx context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}

func (r *karpenterReceiver) scrapeLoop(ctx context.Context) {
	ticker := time.NewTicker(r.cfg.ScrapeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.scrapeMetrics(ctx); err != nil {
				r.logger.Error("Failed to scrape Karpenter metrics", zap.Error(err))
			}
		}
	}
}

func (r *karpenterReceiver) scrapeMetrics(ctx context.Context) error {
	// Simple HTTP GET to scrape metrics
	resp, err := http.Get(r.cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to scrape Karpenter metrics: %w", err)
	}
	defer resp.Body.Close()

	// Create empty metrics for now - in a real implementation you would parse the response
	metrics := pmetric.NewMetrics()
	
	return r.consumer.ConsumeMetrics(ctx, metrics)
}