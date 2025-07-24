package kedareceiver

import (
    "context"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    
    "github.com/prometheus/common/expfmt"
    "github.com/prometheus/common/model"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/consumer/pdata"
    "go.uber.org/zap"
)

type kedaReceiver struct {
    cfg      *Config
    consumer consumer.Metrics
    logger   *zap.Logger
    stopCh   chan struct{}
}

func (r *kedaReceiver) Start(ctx context.Context, host component.Host) error {
    r.stopCh = make(chan struct{})
    go r.scrapeLoop(ctx)
    r.logger.Info("KEDA receiver started")
    return nil
}

func (r *kedaReceiver) scrapeLoop(ctx context.Context) {
    ticker := time.NewTicker(r.cfg.ScrapeInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := r.scrapeAndExport(ctx); err != nil {
                r.logger.Error("Error scraping KEDA metrics", zap.Error(err))
            }
        case <-r.stopCh:
            return
        }
    }
}

func (r *kedaReceiver) scrapeAndExport(ctx context.Context) error {
    url := fmt.Sprintf("http://%s/metrics", r.cfg.Endpoint)
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to fetch metrics: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
    }
    
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }
    
    parser := expfmt.TextParser{}
    metricFamilies, err := parser.TextToMetricFamilies(string(data))
    if err != nil {
        return fmt.Errorf("failed to parse metrics: %w", err)
    }
    
    metrics := pdata.NewMetrics()
    rm := metrics.ResourceMetrics().AppendEmpty()
    sm := rm.ScopeMetrics().AppendEmpty()
    
    for name, mf := range metricFamilies {
        m := sm.Metrics().AppendEmpty()
        m.SetName(name)
        
        // Convert Prometheus metrics to OpenTelemetry metrics
        switch mf.GetType() {
        case dto.MetricType_GAUGE:
            dps := m.SetEmptyGauge().DataPoints()
            for _, metric := range mf.GetMetric() {
                dp := dps.AppendEmpty()
                dp.SetTimestamp(pdata.NewTimestampFromTime(time.Now()))
                dp.SetDoubleVal(metric.GetGauge().GetValue())
                
                // Add labels
                for _, label := range metric.GetLabel() {
                    dp.Attributes().InsertString(label.GetName(), label.GetValue())
                }
            }
        case dto.MetricType_COUNTER:
            dps := m.SetEmptySum().DataPoints()
            m.Sum().SetIsMonotonic(true)
            m.Sum().SetAggregationTemporality(pdata.AggregationTemporalityCumulative)
            
            for _, metric := range mf.GetMetric() {
                dp := dps.AppendEmpty()
                dp.SetTimestamp(pdata.NewTimestampFromTime(time.Now()))
                dp.SetDoubleVal(metric.GetCounter().GetValue())
                
                // Add labels
                for _, label := range metric.GetLabel() {
                    dp.Attributes().InsertString(label.GetName(), label.GetValue())
                }
            }
        case dto.MetricType_HISTOGRAM:
            dps := m.SetEmptyHistogram().DataPoints()
            m.Histogram().SetAggregationTemporality(pdata.AggregationTemporalityCumulative)
            
            for _, metric := range mf.GetMetric() {
                dp := dps.AppendEmpty()
                dp.SetTimestamp(pdata.NewTimestampFromTime(time.Now()))
                
                hist := metric.GetHistogram()
                dp.SetCount(uint64(hist.GetSampleCount()))
                dp.SetSum(hist.GetSampleSum())
                
                // Add buckets
                for i, bucket := range hist.GetBucket() {
                    if i >= len(hist.GetBucket()) {
                        break
                    }
                    
                    dp.BucketCounts().Append(uint64(bucket.GetCumulativeCount()))
                    if i < len(hist.GetBucket())-1 {
                        dp.ExplicitBounds().Append(hist.GetBucket()[i+1].GetUpperBound())
                    }
                }
                
                // Add labels
                for _, label := range metric.GetLabel() {
                    dp.Attributes().InsertString(label.GetName(), label.GetValue())
                }
            }
        }
    }
    
    err = r.consumer.ConsumeMetrics(ctx, metrics)
    if err != nil {
        return fmt.Errorf("failed to consume metrics: %w", err)
    }
    
    return nil
}

func (r *kedaReceiver) Shutdown(ctx context.Context) error {
    close(r.stopCh)
    r.logger.Info("KEDA receiver stopped")
    return nil
}

func createMetricsReceiver(
    ctx context.Context,
    settings component.ReceiverCreateSettings,
    cfg component.Config,
    nextConsumer consumer.Metrics,
) (component.MetricsReceiver, error) {
    rcfg := cfg.(*Config)
    return &kedaReceiver{
        cfg:      rcfg,
        consumer: nextConsumer,
        logger:   settings.Logger,
    }, nil
}
