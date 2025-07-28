package karpenterreceiver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/common/expfmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type karpenterReceiver struct {
	cfg      *Config
	consumer consumer.Metrics
	logger   *zap.Logger
	ticker   *time.Ticker
	stopCh   chan struct{}
}

func newKarpenterReceiver(cfg *Config, consumer consumer.Metrics, logger *zap.Logger) *karpenterReceiver {
	return &karpenterReceiver{
		cfg:      cfg,
		consumer: consumer,
		logger:   logger,
		ticker:   time.NewTicker(cfg.ScrapeInterval),
		stopCh:   make(chan struct{}),
	}
}

func (r *karpenterReceiver) Start(ctx context.Context, host component.Host) error {
	go r.scrapeLoop(ctx)
	r.logger.Info("Karpenter receiver started")
	return nil
}

func (r *karpenterReceiver) scrapeLoop(ctx context.Context) {
	for {
		select {
		case <-r.ticker.C:
			if err := r.scrapeAndSend(ctx); err != nil {
				r.logger.Error("Error scraping Karpenter metrics", zap.Error(err))
			}
		case <-r.stopCh:
			return
		}
	}
}

func (r *karpenterReceiver) scrapeAndSend(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/metrics", r.cfg.Endpoint)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	parser := expfmt.TextParser{}
	families, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return err
	}

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	sm.Scope().SetName("karpenter")

	now := pmetric.NewTimestampFromTime(time.Now())
	for name, mf := range families {
		for _, m := range mf.Metric {
			metric := sm.Metrics().AppendEmpty()
			metric.SetName(name)
			metric.SetDataType(pmetric.MetricTypeGauge)
			dp := metric.Gauge().DataPoints().AppendEmpty()
			dp.SetTimestamp(now)
			if m.Gauge != nil {
				dp.SetDoubleValue(m.Gauge.GetValue())
			}
			for _, l := range m.Label {
				dp.Attributes().PutStr(l.GetName(), l.GetValue())
			}
		}
	}

	return r.consumer.ConsumeMetrics(ctx, md)
}

func (r *karpenterReceiver) Shutdown(context.Context) error {
	r.ticker.Stop()
	close(r.stopCh)
	r.logger.Info("Karpenter receiver stopped")
	return nil
}