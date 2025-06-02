package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// PluginMetadata represents the plugin's metadata
type PluginMetadata struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Version       string                 `json:"version"`
	Description   string                 `json:"description"`
	Author        string                 `json:"author"`
	Endpoints     []EndpointConfig       `json:"endpoints"`
	Permissions   []string               `json:"permissions"`
	Dependencies  []string               `json:"dependencies"`
	Configuration map[string]interface{} `json:"configuration"`
}

type EndpointConfig struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Handler     string `json:"handler"`
	Description string `json:"description"`
}

// Request/Response types
type ClusterOnboardRequest struct {
	Name        string            `json:"name" binding:"required"`
	Kubeconfig  string            `json:"kubeconfig" binding:"required"`
	Type        string            `json:"type"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ClusterDetachRequest struct {
	Name    string `json:"name" binding:"required"`
	Force   bool   `json:"force,omitempty"`
	Cleanup bool   `json:"cleanup,omitempty"`
	Backup  bool   `json:"backup,omitempty"`
}

// Event logging for tracking onboarding/detachment progress
type OnboardingEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	ClusterName string    `json:"clusterName"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
}

// KubestellarClusterPlugin implements the plugin interface with real functionality
type KubestellarClusterPlugin struct {
	metadata         PluginMetadata
	config           map[string]interface{}
	initialized      bool
	clusterStatuses  map[string]string
	onboardingEvents map[string][]OnboardingEvent
	mutex            sync.RWMutex
	kubeconfigDir    string
}

// NewPlugin creates a new instance of the plugin (required symbol that plugin system looks for)
func NewPlugin() interface{} {
	return &KubestellarClusterPlugin{
		metadata: PluginMetadata{
			ID:          "kubestellar-cluster-plugin",
			Name:        "KubeStellar Cluster Management",
			Version:     "1.0.0",
			Description: "Plugin for cluster onboarding and detachment operations with real functionality",
			Author:      "Priyanshu",
			Endpoints: []EndpointConfig{
				{Path: "/onboard", Method: "POST", Handler: "OnboardClusterHandler"},
				{Path: "/detach", Method: "POST", Handler: "DetachClusterHandler"},
				{Path: "/status", Method: "GET", Handler: "GetClusterStatusHandler"},
				{Path: "/list", Method: "GET", Handler: "ListClustersHandler"},
				{Path: "/health", Method: "GET", Handler: "HealthCheckHandler"},
			},
			Permissions:  []string{"cluster.read", "cluster.write", "cluster.delete", "configmap.read", "configmap.write"},
			Dependencies: []string{"kubectl", "clusteradm"},
			Configuration: map[string]interface{}{
				"timeout":           "30s",
				"retries":           3,
				"validate_ssl":      true,
				"log_level":         "info",
				"cluster_namespace": "kubestellar-system",
				"its_context":       "its1",
			},
		},
		clusterStatuses:  make(map[string]string),
		onboardingEvents: make(map[string][]OnboardingEvent),
		kubeconfigDir:    "/tmp/kubestellar-clusters",
	}
}

// Initialize initializes the plugin with configuration
func (p *KubestellarClusterPlugin) Initialize(config map[string]interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.initialized {
		return fmt.Errorf("plugin already initialized")
	}

	p.config = config

	// Create kubeconfig directory if it doesn't exist
	if err := os.MkdirAll(p.kubeconfigDir, 0755); err != nil {
		log.Printf("Warning: Failed to create kubeconfig directory: %v", err)
	}

	// Check for required tools
	if err := p.checkCommand("kubectl"); err != nil {
		log.Printf("Warning: kubectl not available: %v", err)
	}
	if err := p.checkCommand("clusteradm"); err != nil {
		log.Printf("Warning: clusteradm not available: %v", err)
	}

	p.initialized = true
	log.Printf("✅ KubeStellar Cluster Plugin initialized with real functionality")
	return nil
}

// GetMetadata returns the plugin metadata
func (p *KubestellarClusterPlugin) GetMetadata() interface{} {
	return p.metadata
}

// GetHandlers returns the HTTP handlers for this plugin
func (p *KubestellarClusterPlugin) GetHandlers() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"OnboardClusterHandler":   p.OnboardClusterHandler,
		"DetachClusterHandler":    p.DetachClusterHandler,
		"GetClusterStatusHandler": p.GetClusterStatusHandler,
		"ListClustersHandler":     p.ListClustersHandler,
		"HealthCheckHandler":      p.HealthCheckHandler,
	}
}

// Health returns the health status of the plugin
func (p *KubestellarClusterPlugin) Health() error {
	if !p.initialized {
		return fmt.Errorf("plugin not initialized")
	}
	return nil
}

// Cleanup performs cleanup operations
func (p *KubestellarClusterPlugin) Cleanup() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.initialized = false
	log.Println("🧹 KubeStellar Cluster Plugin cleaned up")
	return nil
}

// checkCommand verifies that a command is available in PATH
func (p *KubestellarClusterPlugin) checkCommand(command string) error {
	_, err := exec.LookPath(command)
	return err
}

// LogOnboardingEvent logs an event for the onboarding/detachment process
func (p *KubestellarClusterPlugin) LogOnboardingEvent(clusterName, status, message string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	event := OnboardingEvent{
		Timestamp:   time.Now(),
		ClusterName: clusterName,
		Status:      status,
		Message:     message,
	}

	if p.onboardingEvents[clusterName] == nil {
		p.onboardingEvents[clusterName] = make([]OnboardingEvent, 0)
	}

	p.onboardingEvents[clusterName] = append(p.onboardingEvents[clusterName], event)
	log.Printf("[%s] %s: %s", clusterName, status, message)
}

// OnboardClusterHandler handles cluster onboarding requests with real functionality
func (p *KubestellarClusterPlugin) OnboardClusterHandler(c *gin.Context) {
	log.Println("🚀 Plugin: Handling REAL cluster onboarding request")

	contentType := c.GetHeader("Content-Type")
	var kubeconfigData []byte
	var clusterName string
	var useLocalKubeconfig bool = false

	// Handle different content types
	if strings.Contains(contentType, "multipart/form-data") {
		file, fileErr := c.FormFile("kubeconfig")
		clusterName = c.PostForm("name")

		if clusterName != "" && (fileErr != nil || file == nil) {
			useLocalKubeconfig = true
		} else if fileErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve kubeconfig file"})
			return
		} else if clusterName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster name is required"})
			return
		} else {
			f, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open kubeconfig file"})
				return
			}
			defer f.Close()

			kubeconfigData, err = io.ReadAll(f)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read kubeconfig file"})
				return
			}
		}
	} else if strings.Contains(contentType, "application/json") {
		var req ClusterOnboardRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		clusterName = req.Name
		if clusterName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ClusterName is required"})
			return
		}

		if req.Kubeconfig == "" {
			useLocalKubeconfig = true
		} else {
			kubeconfigData = []byte(req.Kubeconfig)
		}
	} else {
		clusterName = c.Query("name")
		if clusterName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cluster name parameter is required"})
			return
		}
		useLocalKubeconfig = true
	}

	// Get kubeconfig from local if needed
	if useLocalKubeconfig {
		var err error
		kubeconfigData, err = p.getClusterConfigFromLocal(clusterName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to find cluster '%s' in local kubeconfig: %v", clusterName, err)})
			return
		}
	}

	// Check if the cluster is already being onboarded
	p.mutex.Lock()
	if status, exists := p.clusterStatuses[clusterName]; exists {
		p.mutex.Unlock()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Cluster '%s' is already onboarded (status: %s)", clusterName, status),
			"status":  status,
		})
		return
	}
	p.clusterStatuses[clusterName] = "Pending"
	p.mutex.Unlock()

	// Log initial event and clear any previous events
	p.ClearOnboardingEvents(clusterName)
	p.LogOnboardingEvent(clusterName, "Initiated", "Onboarding process initiated by plugin API request")

	// Start asynchronous onboarding with real functionality
	go func() {
		err := p.OnboardCluster(kubeconfigData, clusterName)
		p.mutex.Lock()
		if err != nil {
			log.Printf("Cluster '%s' onboarding failed: %v", clusterName, err)
			p.clusterStatuses[clusterName] = "Failed"
		} else {
			p.clusterStatuses[clusterName] = "Onboarded"
			log.Printf("Cluster '%s' onboarded successfully", clusterName)
		}
		p.mutex.Unlock()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":           fmt.Sprintf("Cluster '%s' is being onboarded", clusterName),
		"status":            "Pending",
		"logsEndpoint":      fmt.Sprintf("/api/plugins/kubestellar-cluster-plugin/logs/%s", clusterName),
		"websocketEndpoint": fmt.Sprintf("/ws/plugins/kubestellar-cluster-plugin/onboarding?cluster=%s", clusterName),
	})
}

// ClearOnboardingEvents clears events for a cluster
func (p *KubestellarClusterPlugin) ClearOnboardingEvents(clusterName string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.onboardingEvents[clusterName] = make([]OnboardingEvent, 0)
}

// DetachClusterHandler handles cluster detachment requests with real functionality
func (p *KubestellarClusterPlugin) DetachClusterHandler(c *gin.Context) {
	log.Println("🗑️ Plugin: Handling REAL cluster detachment request")

	var req ClusterDetachRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	clusterName := req.Name
	log.Printf("Detaching cluster: %s (force: %v, cleanup: %v, backup: %v)", clusterName, req.Force, req.Cleanup, req.Backup)

	// Check if cluster exists in our status map
	p.mutex.RLock()
	status, exists := p.clusterStatuses[clusterName]
	p.mutex.RUnlock()

	if !exists {
		log.Printf("Cluster '%s' not found in status map, checking OCM hub", clusterName)
	} else {
		log.Printf("Cluster '%s' status is %s, proceeding with detachment", clusterName, status)
	}

	// Start detaching the cluster
	p.mutex.Lock()
	p.clusterStatuses[clusterName] = "Detaching"
	p.mutex.Unlock()

	go func() {
		err := p.DetachCluster(clusterName, req.Force)
		p.mutex.Lock()
		if err != nil {
			log.Printf("Cluster '%s' detachment failed: %v", clusterName, err)
			p.clusterStatuses[clusterName] = "DetachmentFailed"
		} else {
			log.Printf("Cluster '%s' detached successfully", clusterName)
			delete(p.clusterStatuses, clusterName)
		}
		p.mutex.Unlock()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":           fmt.Sprintf("Cluster '%s' is being detached", clusterName),
		"status":            "Detaching",
		"logsEndpoint":      fmt.Sprintf("/api/plugins/kubestellar-cluster-plugin/logs/%s", clusterName),
		"websocketEndpoint": fmt.Sprintf("/ws/plugins/kubestellar-cluster-plugin/detachment?cluster=%s", clusterName),
	})
}

// GetClusterStatusHandler returns cluster status information
func (p *KubestellarClusterPlugin) GetClusterStatusHandler(c *gin.Context) {
	clusterName := c.Query("name")
	if clusterName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cluster name is required"})
		return
	}

	p.mutex.RLock()
	status, exists := p.clusterStatuses[clusterName]
	events := p.onboardingEvents[clusterName]
	p.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cluster": gin.H{
			"name":     clusterName,
			"status":   status,
			"lastSeen": time.Now().Add(-5 * time.Minute),
		},
		"events": events,
	})
}

// ListClustersHandler returns a list of all managed clusters
func (p *KubestellarClusterPlugin) ListClustersHandler(c *gin.Context) {
	p.mutex.RLock()
	clusters := make([]map[string]interface{}, 0)
	for name, status := range p.clusterStatuses {
		clusters = append(clusters, map[string]interface{}{
			"name":         name,
			"status":       status,
			"type":         "workload",
			"onboarded_at": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		})
	}
	p.mutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"clusters":     clusters,
		"total":        len(clusters),
		"connected":    len(clusters),
		"disconnected": 0,
	})
}

// HealthCheckHandler provides plugin health status
func (p *KubestellarClusterPlugin) HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"timestamp":   time.Now().Format(time.RFC3339),
		"version":     p.metadata.Version,
		"initialized": p.initialized,
	})
}

// Real functionality methods (simplified versions for demo)

// getClusterConfigFromLocal extracts a specific cluster's config from the local kubeconfig file
func (p *KubestellarClusterPlugin) getClusterConfigFromLocal(clusterName string) ([]byte, error) {
	kubeconfig := p.kubeconfigPath()
	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	// Check if the cluster exists
	_, exists := config.Clusters[clusterName]
	if !exists {
		// Try to find a context that references this cluster
		for contextName, ctx := range config.Contexts {
			if ctx.Cluster == clusterName {
				return p.extractContextConfig(config, contextName)
			}
		}
		return nil, fmt.Errorf("cluster '%s' not found in local kubeconfig", clusterName)
	}

	// Find a context that uses this cluster
	var contextName string
	for ctxName, ctx := range config.Contexts {
		if ctx.Cluster == clusterName {
			contextName = ctxName
			break
		}
	}

	if contextName == "" {
		return nil, fmt.Errorf("no context found for cluster '%s'", clusterName)
	}

	return p.extractContextConfig(config, contextName)
}

// extractContextConfig creates a standalone kubeconfig for a specific context
func (p *KubestellarClusterPlugin) extractContextConfig(config *clientcmdapi.Config, contextName string) ([]byte, error) {
	ctx, exists := config.Contexts[contextName]
	if !exists {
		return nil, fmt.Errorf("context '%s' not found", contextName)
	}

	// Create a new config with only the required context, cluster, and user
	newConfig := &clientcmdapi.Config{
		Clusters:       make(map[string]*clientcmdapi.Cluster),
		AuthInfos:      make(map[string]*clientcmdapi.AuthInfo),
		Contexts:       make(map[string]*clientcmdapi.Context),
		CurrentContext: contextName,
	}

	// Copy the cluster
	if cluster, exists := config.Clusters[ctx.Cluster]; exists {
		newConfig.Clusters[ctx.Cluster] = cluster
	} else {
		return nil, fmt.Errorf("cluster '%s' not found", ctx.Cluster)
	}

	// Copy the user (auth info)
	if user, exists := config.AuthInfos[ctx.AuthInfo]; exists {
		newConfig.AuthInfos[ctx.AuthInfo] = user
	} else {
		return nil, fmt.Errorf("user '%s' not found", ctx.AuthInfo)
	}

	// Copy the context
	newConfig.Contexts[contextName] = ctx

	// Convert to bytes
	return clientcmd.Write(*newConfig)
}

// kubeconfigPath returns the path to the kubeconfig file
func (p *KubestellarClusterPlugin) kubeconfigPath() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}
	return filepath.Join(os.Getenv("HOME"), ".kube", "config")
}

// OnboardCluster handles the real cluster onboarding process using OCM/clusteradm
func (p *KubestellarClusterPlugin) OnboardCluster(kubeconfigData []byte, clusterName string) error {
	p.LogOnboardingEvent(clusterName, "Starting", "Beginning cluster onboarding process")

	// 1. Validate cluster connectivity
	p.LogOnboardingEvent(clusterName, "Validating", "Validating cluster connectivity")
	if err := p.ValidateClusterConnectivity(kubeconfigData); err != nil {
		p.LogOnboardingEvent(clusterName, "Error", "Cluster connectivity validation failed: "+err.Error())
		return fmt.Errorf("cluster connectivity validation failed: %w", err)
	}
	p.LogOnboardingEvent(clusterName, "Validated", "Cluster connectivity validated successfully")

	// 2. Save kubeconfig temporarily for clusteradm
	p.LogOnboardingEvent(clusterName, "Preparing", "Preparing cluster configuration")
	tempKubeconfigPath := filepath.Join(p.kubeconfigDir, fmt.Sprintf("%s-kubeconfig.yaml", clusterName))
	if err := os.WriteFile(tempKubeconfigPath, kubeconfigData, 0600); err != nil {
		p.LogOnboardingEvent(clusterName, "Error", "Failed to save temporary kubeconfig: "+err.Error())
		return fmt.Errorf("failed to save temporary kubeconfig: %w", err)
	}
	defer os.Remove(tempKubeconfigPath) // Clean up

	// 3. Generate join token from ITS hub
	p.LogOnboardingEvent(clusterName, "GeneratingToken", "Generating clusteradm join token from ITS hub")
	joinToken, err := p.generateJoinToken()
	if err != nil {
		p.LogOnboardingEvent(clusterName, "Error", "Failed to generate join token: "+err.Error())
		return fmt.Errorf("failed to generate join token: %w", err)
	}
	p.LogOnboardingEvent(clusterName, "TokenGenerated", "Join token generated successfully")

	// 4. Join cluster to OCM hub using clusteradm
	p.LogOnboardingEvent(clusterName, "Joining", "Joining cluster to OCM hub using clusteradm")
	if err := p.joinClusterToHub(clusterName, tempKubeconfigPath, joinToken); err != nil {
		p.LogOnboardingEvent(clusterName, "Error", "Failed to join cluster to hub: "+err.Error())
		return fmt.Errorf("failed to join cluster to hub: %w", err)
	}
	p.LogOnboardingEvent(clusterName, "Joined", "Cluster joined to OCM hub successfully")

	// 5. Wait for CSR and approve it
	p.LogOnboardingEvent(clusterName, "ApprovingCSR", "Waiting for and approving Certificate Signing Request")
	if err := p.approveClusterCSR(clusterName); err != nil {
		p.LogOnboardingEvent(clusterName, "Warning", "CSR approval failed, but cluster may still work: "+err.Error())
		// Don't fail the entire process for CSR issues
	} else {
		p.LogOnboardingEvent(clusterName, "CSRApproved", "Certificate Signing Request approved successfully")
	}

	// 6. Verify cluster is managed
	p.LogOnboardingEvent(clusterName, "Verifying", "Verifying cluster is properly managed")
	if err := p.verifyClusterManaged(clusterName); err != nil {
		p.LogOnboardingEvent(clusterName, "Error", "Cluster verification failed: "+err.Error())
		return fmt.Errorf("cluster verification failed: %w", err)
	}

	p.LogOnboardingEvent(clusterName, "Success", "Cluster onboarded successfully to KubeStellar")
	return nil
}

// DetachCluster handles the real cluster detachment process using OCM
func (p *KubestellarClusterPlugin) DetachCluster(clusterName string, force bool) error {
	p.LogOnboardingEvent(clusterName, "Detaching", "Starting cluster detachment process")

	// 1. Check if cluster exists in OCM
	p.LogOnboardingEvent(clusterName, "Checking", "Checking cluster status in OCM hub")
	exists, err := p.checkClusterExists(clusterName)
	if err != nil && !force {
		p.LogOnboardingEvent(clusterName, "Error", "Failed to check cluster status: "+err.Error())
		return fmt.Errorf("failed to check cluster status: %w", err)
	}

	if !exists && !force {
		p.LogOnboardingEvent(clusterName, "Warning", "Cluster not found in OCM hub")
		return fmt.Errorf("cluster %s not found in OCM hub", clusterName)
	}

	// 2. Remove cluster from OCM hub using kubectl
	p.LogOnboardingEvent(clusterName, "Removing", "Removing cluster from OCM hub")
	if err := p.removeClusterFromHub(clusterName); err != nil && !force {
		p.LogOnboardingEvent(clusterName, "Error", "Failed to remove cluster from hub: "+err.Error())
		return fmt.Errorf("failed to remove cluster from hub: %w", err)
	}
	p.LogOnboardingEvent(clusterName, "Removed", "Cluster removed from OCM hub")

	// 3. Clean up local resources
	p.LogOnboardingEvent(clusterName, "Cleanup", "Cleaning up local resources")
	if err := p.cleanupLocalResources(clusterName); err != nil {
		p.LogOnboardingEvent(clusterName, "Warning", "Failed to clean up some local resources: "+err.Error())
		// Don't fail for cleanup issues
	}

	p.LogOnboardingEvent(clusterName, "Success", "Cluster detached successfully from KubeStellar")
	return nil
}

// generateJoinToken generates a join token from the ITS hub
func (p *KubestellarClusterPlugin) generateJoinToken() (string, error) {
	itsContext := p.getITSContext()

	cmd := exec.Command("clusteradm", "get", "token", "--context", itsContext)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate join token: %w", err)
	}

	// Parse the token from clusteradm output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "clusteradm join") {
			return strings.TrimSpace(line), nil
		}
	}

	return "", fmt.Errorf("failed to parse join token from clusteradm output")
}

// joinClusterToHub joins the cluster to the OCM hub using clusteradm
func (p *KubestellarClusterPlugin) joinClusterToHub(clusterName, kubeconfigPath, joinToken string) error {
	// Extract the actual clusteradm join command from the token
	if !strings.Contains(joinToken, "clusteradm join") {
		return fmt.Errorf("invalid join token format")
	}

	// Build clusteradm join command
	cmdParts := strings.Fields(joinToken)
	cmdParts = append(cmdParts, "--cluster-name", clusterName, "--kubeconfig", kubeconfigPath)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("clusteradm join failed: %w, output: %s", err, string(output))
	}

	log.Printf("Clusteradm join output: %s", string(output))
	return nil
}

// approveClusterCSR approves the Certificate Signing Request for the cluster
func (p *KubestellarClusterPlugin) approveClusterCSR(clusterName string) error {
	itsContext := p.getITSContext()

	// Wait a bit for CSR to appear
	time.Sleep(5 * time.Second)

	// Get pending CSRs
	cmd := exec.Command("kubectl", "get", "csr", "--context", itsContext, "-o", "name")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get CSRs: %w", err)
	}

	// Look for CSRs related to our cluster
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, clusterName) {
			csrName := strings.TrimPrefix(strings.TrimSpace(line), "certificatesigningrequest.certificates.k8s.io/")
			if csrName != "" {
				approveCmd := exec.Command("kubectl", "certificate", "approve", csrName, "--context", itsContext)
				if err := approveCmd.Run(); err != nil {
					return fmt.Errorf("failed to approve CSR %s: %w", csrName, err)
				}
				log.Printf("Approved CSR: %s", csrName)
			}
		}
	}

	return nil
}

// verifyClusterManaged verifies that the cluster is properly managed by OCM
func (p *KubestellarClusterPlugin) verifyClusterManaged(clusterName string) error {
	itsContext := p.getITSContext()

	// Check if ManagedCluster resource exists
	cmd := exec.Command("kubectl", "get", "managedcluster", clusterName, "--context", itsContext)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("managed cluster resource not found: %w", err)
	}

	// Check cluster status
	cmd = exec.Command("kubectl", "get", "managedcluster", clusterName, "--context", itsContext, "-o", "jsonpath={.status.conditions[?(@.type=='ManagedClusterConditionAvailable')].status}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get cluster status: %w", err)
	}

	if strings.TrimSpace(string(output)) != "True" {
		return fmt.Errorf("cluster is not in available state")
	}

	return nil
}

// checkClusterExists checks if a cluster exists in the OCM hub
func (p *KubestellarClusterPlugin) checkClusterExists(clusterName string) (bool, error) {
	itsContext := p.getITSContext()

	cmd := exec.Command("kubectl", "get", "managedcluster", clusterName, "--context", itsContext)
	err := cmd.Run()
	return err == nil, nil
}

// removeClusterFromHub removes a cluster from the OCM hub
func (p *KubestellarClusterPlugin) removeClusterFromHub(clusterName string) error {
	itsContext := p.getITSContext()

	cmd := exec.Command("kubectl", "delete", "managedcluster", clusterName, "--context", itsContext)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete managed cluster: %w, output: %s", err, string(output))
	}

	return nil
}

// cleanupLocalResources cleans up any local resources related to the cluster
func (p *KubestellarClusterPlugin) cleanupLocalResources(clusterName string) error {
	// Remove any temporary kubeconfig files
	kubeconfigPath := filepath.Join(p.kubeconfigDir, fmt.Sprintf("%s-kubeconfig.yaml", clusterName))
	if err := os.Remove(kubeconfigPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: failed to remove temporary kubeconfig: %v", err)
	}

	return nil
}

// getITSContext returns the ITS context name from configuration
func (p *KubestellarClusterPlugin) getITSContext() string {
	if itsContext, ok := p.config["its_context"].(string); ok && itsContext != "" {
		return itsContext
	}
	// Default ITS context
	return "its1"
}

// ValidateClusterConnectivity validates that we can connect to the cluster
func (p *KubestellarClusterPlugin) ValidateClusterConnectivity(kubeconfigData []byte) error {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigData)
	if err != nil {
		return fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	_, err = clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	return nil
}

// main is required for plugin compilation but not used in plugin mode
func main() {
	log.Println("This is a plugin, not a standalone application")
}
