package nodemanagement

import (
	"digitalcashtools/monerod-proxy/httpclient"
	"fmt"

	"gopkg.in/ini.v1"
)

type INodeProvider interface {
	GetBaseURL() string
	GetAnyNodesAvailable() bool
	ReportNodeConnectionFailure()
	CheckNodeHealth()
}

type NodeInfo struct {
	URL             string
	PassedLastCheck bool
}

type NodeProvider struct {
	SelectedNodeIndex int
	Nodes             []NodeInfo
	AnyNodesAvailable bool
}

func (nodeProvider *NodeProvider) GetBaseURL() string {
	return nodeProvider.Nodes[nodeProvider.SelectedNodeIndex].URL
}

func (nodeProvider *NodeProvider) GetAnyNodesAvailable() bool {
	return nodeProvider.AnyNodesAvailable
}

func (nodeProvider *NodeProvider) CheckNodeHealth() {
	for i := 0; i < len(nodeProvider.Nodes); i++ {
		_, statusCode, err := httpclient.ExecuteGETRequest(nodeProvider.Nodes[i].URL + "get_height")

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
		if nodeProvider.Nodes[i].PassedLastCheck {
			chosenNodeIndex = i
			availableNodeFound = true
			break
		}
	}

	nodeProvider.SelectedNodeIndex = chosenNodeIndex
	nodeProvider.AnyNodesAvailable = availableNodeFound
}

func LoadNodeProviderFromConfig(cfg *ini.File) *NodeProvider {
	nodes := []NodeInfo{}
	nodeURLs := cfg.Section("").Key("node").Strings(",")
	fmt.Println(nodeURLs)

	for i := 0; i < len(nodeURLs); i++ {
		fmt.Println("Adding URL: " + nodeURLs[i])

		nodeInfo := NodeInfo{
			URL:             nodeURLs[i],
			PassedLastCheck: true,
		}

		nodes = append(nodes, nodeInfo)
	}

	nodeProvider := &NodeProvider{
		SelectedNodeIndex: 0,
		Nodes:             nodes,
		AnyNodesAvailable: false,
	}

	return nodeProvider
}

func (nodeProvider *NodeProvider) ReportNodeConnectionFailure() {
	fmt.Println("Detected node failure:\t", nodeProvider.GetBaseURL())

	nodeProvider.CheckNodeHealth()

	if nodeProvider.GetAnyNodesAvailable() {
		fmt.Println("Switched to node:\t", nodeProvider.GetBaseURL())
	} else {
		fmt.Println("All nodes failed health check, no nodes available.")
	}
}
