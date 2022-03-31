package nodemanagement

import (
	"gopkg.in/ini.v1"
)

type INodeProvider interface {
	GetBaseURL() string
}

type NodeInfo struct {
	URL string
}

type NodeProvider struct {
	SelectedNodeIndex int
	Nodes             []NodeInfo
}

func (nodeProvider *NodeProvider) GetBaseURL() string {
	return nodeProvider.Nodes[nodeProvider.SelectedNodeIndex].URL
}

// func checkNodeHealth(nodeProvider NodeProvider) {
//    Ping each node in nodeProvider.Nodes and set SelectedNodeIndex to the one with the best ping time.
// }

func LoadNodeProviderFromConfig(cfg *ini.File) *NodeProvider {
	nodeInfoSlice := []NodeInfo{}
	baseURL := cfg.Section("").Key("node").Value()
	nodeInfo := NodeInfo{
		URL: baseURL,
	}

	nodeInfoSlice = append(nodeInfoSlice, nodeInfo)

	nodeProvider := &NodeProvider{
		SelectedNodeIndex: 0,
		Nodes:             nodeInfoSlice,
	}

	return nodeProvider
}
