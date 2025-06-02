# KubeStellar Cluster Management Plugin

A dynamic plugin for KubeStellar that provides cluster onboarding and detachment functionality with real OCM integration.

## 🚀 Features

- **Real Cluster Onboarding**: Add clusters to KubeStellar with kubeconfig validation and OCM integration
- **Safe Cluster Detachment**: Remove clusters with proper cleanup and confirmation
- **Progress Tracking**: Real-time onboarding/detachment progress with detailed event logging
- **Multi-format Support**: Accept kubeconfig via file upload, JSON, or local context reference
- **Health Monitoring**: Built-in health checks and status monitoring

## 📦 Installation

### Via GitHub URL (Recommended)

1. Open your KubeStellar UI
2. Navigate to **Plugins** → **Plugin Store**
3. Click **"Install from GitHub"**
4. Enter this repository URL:
   ```
   https://github.com/your-username/kubestellar-cluster-plugin
   ```
5. Click **"Install"**

**That's it!** The plugin system will automatically:
- 📥 Clone this repository
- 🔨 Build the plugin using `go build -buildmode=plugin`
- 🚀 Load it into your KubeStellar system
- ✅ Register all API endpoints

### Manual Installation (Alternative)

If you prefer to install locally:

1. Clone this repository
2. Navigate to **Plugins** → **Plugin Store** → **"Install from File"**
3. Upload both `plugin.yaml` and `main.go` files

## 🎯 Quick Start

Once installed, the plugin provides these endpoints:

### Onboard a Cluster (with kubeconfig upload)
```bash
curl -X POST "http://localhost:8080/api/plugins/kubestellar-cluster-plugin/onboard" \
  -F "name=my-cluster" \
  -F "kubeconfig=@/path/to/kubeconfig.yaml"
```

### Onboard from Local Context
```bash
curl -X POST "http://localhost:8080/api/plugins/kubestellar-cluster-plugin/onboard?name=existing-context"
```

### Detach a Cluster
```bash
curl -X POST "http://localhost:8080/api/plugins/kubestellar-cluster-plugin/detach" \
  -H "Content-Type: application/json" \
  -d '{"name": "my-cluster", "force": false, "cleanup": true}'
```

### Monitor Status
```bash
curl "http://localhost:8080/api/plugins/kubestellar-cluster-plugin/status?name=my-cluster"
```

## 🏗️ How It Works

This plugin integrates with KubeStellar's Open Cluster Management (OCM) system:

1. **Validates** kubeconfig connectivity to ensure cluster is accessible
2. **Generates** clusteradm join tokens from the ITS hub
3. **Joins** clusters using `clusteradm join` command
4. **Approves** Certificate Signing Requests (CSRs) automatically
5. **Monitors** cluster health and manages lifecycle

## 🛠️ Development

### Prerequisites
- Go 1.21+
- kubectl (Kubernetes CLI)
- clusteradm (OCM CLI tool)
- Access to a KubeStellar environment

### Simple Plugin Structure
This plugin follows the minimal Jenkins-style structure:

```
kubestellar-cluster-plugin/
├── plugin.yaml    # Plugin metadata and API definitions
├── main.go        # Complete plugin implementation
└── README.md      # This documentation
```

**No complex build systems needed!** The KubeStellar plugin system handles compilation automatically.

### Key Functions
- `NewPlugin()` - Plugin factory function (required by plugin system)
- `Initialize()` - Setup and dependency validation
- `GetHandlers()` - Returns HTTP endpoint handlers
- `GetMetadata()` - Plugin information and configuration

## 🔧 Configuration

The plugin automatically configures itself with sensible defaults:

```yaml
timeout: 30s
retries: 3
validate_ssl: true
log_level: info
cluster_namespace: kubestellar-system
its_context: its1
```

## 🤝 Contributing

1. Fork this repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is part of the CNCF KubeStellar ecosystem and follows the project's licensing terms.

## 🙋‍♂️ Support

- 📖 [KubeStellar Documentation](https://docs.kubestellar.io)
- 💬 [KubeStellar Slack](https://kubernetes.slack.com/channels/kubestellar)
- 🐛 [Report Issues](https://github.com/your-username/kubestellar-cluster-plugin/issues)

---

Made with ❤️ for the CNCF LFX Mentorship Program 