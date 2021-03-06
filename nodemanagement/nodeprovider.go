package nodemanagement

import (
	"digitalcashtools/monerod-proxy/httpclient"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type INodeProvider interface {
	GetBaseURL() string
	GetAnyNodesAvailable() bool
	ReportNodeConnectionFailure()
	CheckNodeHealth()
	GetAvailableNodes() []string
	SetNodeEnabled(NodeURL string, enabled bool) bool
}

type NodeInfo struct {
	URL             string
	PassedLastCheck bool
	Enabled         bool
}

type ExecuteGETRequestFunc func(string) (string, int, error)

type NodeProvider struct {
	SelectedNodeIndex     int
	Nodes                 []NodeInfo
	AnyNodesAvailable     bool
	executeGETRequestFunc ExecuteGETRequestFunc
}

func (nodeProvider *NodeProvider) GetBaseURL() string {
	return nodeProvider.Nodes[nodeProvider.SelectedNodeIndex].URL
}

func (nodeProvider *NodeProvider) GetAnyNodesAvailable() bool {
	return nodeProvider.AnyNodesAvailable
}

func (nodeProvider *NodeProvider) GetAvailableNodes() []string {
	result := []string{}

	for i := 0; i < len(nodeProvider.Nodes); i++ {
		if nodeProvider.Nodes[i].PassedLastCheck && nodeProvider.Nodes[i].Enabled {
			result = append(result, nodeProvider.Nodes[i].URL)
		}
	}

	return result
}

func (nodeProvider *NodeProvider) CheckNodeHealth() {
	for i := 0; i < len(nodeProvider.Nodes); i++ {
		if nodeProvider.Nodes[i].PassedLastCheck {
			continue
		}

		if !nodeProvider.Nodes[i].Enabled {
			continue
		}

		_, statusCode, err := nodeProvider.executeGETRequestFunc(nodeProvider.Nodes[i].URL + "get_height")

		if err != nil {
			nodeProvider.Nodes[i].PassedLastCheck = false
			continue
		}

		if !(statusCode >= 200 && statusCode <= 299) {
			nodeProvider.Nodes[i].PassedLastCheck = false
			continue
		}

		nodeProvider.Nodes[i].PassedLastCheck = true
	}

	availableNodeFound := false
	chosenNodeIndex := 0
	for i := 0; i < len(nodeProvider.Nodes); i++ {
		if nodeProvider.Nodes[i].PassedLastCheck && nodeProvider.Nodes[i].Enabled {
			chosenNodeIndex = i
			availableNodeFound = true
			break
		}
	}

	nodeProvider.SelectedNodeIndex = chosenNodeIndex
	nodeProvider.AnyNodesAvailable = availableNodeFound
}

func CreateNodeProviderFromConfig(cfg *ini.File) *NodeProvider {
	nodes := []NodeInfo{}
	nodeURLs := cfg.Section("").Key("node").Strings(",")
	log.Info(nodeURLs)

	nodeProvider := &NodeProvider{
		SelectedNodeIndex:     0,
		Nodes:                 nodes,
		AnyNodesAvailable:     false,
		executeGETRequestFunc: httpclient.ExecuteGETRequest,
	}

	for i := 0; i < len(nodeURLs); i++ {
		log.Debug("Adding URL: " + nodeURLs[i])
		nodeProvider.AddNode(nodeURLs[i])
	}

	return nodeProvider
}

func (nodeProvider *NodeProvider) AddNode(URL string) {
	nodeInfo := NodeInfo{
		URL:             URL,
		PassedLastCheck: true,
		Enabled:         true,
	}

	nodeProvider.Nodes = append(nodeProvider.Nodes, nodeInfo)
}

func (nodeProvider *NodeProvider) ReportNodeConnectionFailure() {
	log.Info("Detected node failure: " + nodeProvider.GetBaseURL())

	nodeProvider.Nodes[nodeProvider.SelectedNodeIndex].PassedLastCheck = false
	nodeProvider.CheckNodeHealth()

	if nodeProvider.GetAnyNodesAvailable() {
		log.Info("Switched to node: " + nodeProvider.GetBaseURL())
	} else {
		log.Info("All nodes failed health check, no nodes available.")
	}
}

func (nodeProvider *NodeProvider) SetNodeEnabled(NodeURL string, enabled bool) bool {
	for i := 0; i < len(nodeProvider.Nodes); i++ {
		if nodeProvider.Nodes[i].URL == NodeURL {
			nodeProvider.Nodes[i].Enabled = enabled
			log.Info("Setting node [", NodeURL, "] to enabled = ", enabled)
			return true
		}
	}

	log.Error("No node found with URL", NodeURL)
	return false
}
