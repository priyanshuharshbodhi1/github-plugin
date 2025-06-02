# KubeStellar Cluster Operations Plugin

A comprehensive cluster management plugin for KubeStellar that provides cluster onboarding, detachment, and monitoring capabilities.

## Features

- **Cluster Onboarding**: Seamlessly onboard new clusters to KubeStellar
- **Cluster Detachment**: Remove clusters from KubeStellar management  
- **Status Monitoring**: Real-time cluster health and status monitoring
- **Event Tracking**: Track cluster onboarding events and operations
- **Self-contained**: No external dependencies on UI backend handlers

## API Endpoints

### POST /onboard
Onboard a new cluster to KubeStellar.

**Request Body:**
```json
{
  "clusterName": "my-cluster",
  "kubeconfig": "base64-encoded-kubeconfig"
}
```

### POST /detach
Detach a cluster from KubeStellar management.

**Request Body:**
```json
{
  "clusterName": "my-cluster"
}
```

### GET /status/:cluster
Get detailed status information for a specific cluster.

### GET /clusters
List all managed clusters with their status information.

### GET /health
Plugin health check endpoint.

### GET /events/:cluster
Get cluster-specific events and operation logs.

## Installation

This plugin can be loaded directly from GitHub:

```bash
# Install from GitHub repository
curl -X POST http://localhost:4000/api/plugins/install \
  -H "Content-Type: application/json" \
  -d '{"source": "https://github.com/priyanshuharshbodhi1/kubestellar-cluster-ops-plugin"}'
```

## Configuration

The plugin supports the following configuration options:

- `timeout`: Operation timeout (default: "60s")
- `cluster_namespace`: Kubernetes namespace for cluster resources (default: "kubestellar-system")
- `its_context`: ITS context name (default: "its1")

## Compatibility

- KubeStellar: >=0.21.0
- Go: >=1.21

## License

This plugin is part of the KubeStellar project and follows the same licensing terms. 