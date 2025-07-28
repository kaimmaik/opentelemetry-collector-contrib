# Karpenter Receiver

The Karpenter receiver collects metrics from Karpenter to monitor node provisioning, scaling decisions, and cluster capacity management.

## Configuration

```yaml
receivers:
  karpenter:
    endpoint: "karpenter.karpenter.svc.cluster.local:8000"
    scrape_interval: 30s
```

### Configuration Parameters

- `endpoint` (default: `karpenter.karpenter.svc.cluster.local:8000`): The Karpenter metrics endpoint
- `scrape_interval` (default: `30s`): How often to scrape metrics

## Metrics

| Name | Description | Unit | Type |
|------|-------------|------|------|
| `karpenter_nodes_total` | Total number of nodes managed by Karpenter | 1 | Gauge |
| `karpenter_pods_pending` | Number of pending pods waiting for resources | 1 | Gauge |
| `karpenter_node_utilization` | Node resource utilization percentage | % | Gauge |

## Example

```yaml
receivers:
  karpenter:
    endpoint: "karpenter.karpenter.svc.cluster.local:8000"
    scrape_interval: 15s

exporters:
  awsemf:
    namespace: EKS/Karpenter
    region: us-west-2

service:
  pipelines:
    metrics:
      receivers: [karpenter]
      exporters: [awsemf]
```