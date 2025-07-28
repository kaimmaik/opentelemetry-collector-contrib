# KEDA Receiver

The KEDA receiver collects metrics from KEDA (Kubernetes Event-driven Autoscaling) operators to monitor autoscaling behavior and performance.

## Configuration

```yaml
receivers:
  keda:
    endpoint: "keda-operator.keda-system.svc.cluster.local:8080"
    scrape_interval: 30s
```

### Configuration Parameters

- `endpoint` (default: `keda-operator.keda-system.svc.cluster.local:8080`): The KEDA operator metrics endpoint
- `scrape_interval` (default: `30s`): How often to scrape metrics

## Metrics

| Name | Description | Unit | Type |
|------|-------------|------|------|
| `keda_scaler_metrics_value` | Current value of the metric used by KEDA for scaling decisions | 1 | Gauge |
| `keda_scaler_active` | Whether the scaler is currently active (1) or inactive (0) | 1 | Gauge |

## Example

```yaml
receivers:
  keda:
    endpoint: "keda-operator.keda-system.svc.cluster.local:8080"
    scrape_interval: 15s

exporters:
  awsemf:
    namespace: EKS/KEDA
    region: us-west-2

service:
  pipelines:
    metrics:
      receivers: [keda]
      exporters: [awsemf]
```