package nodemanagement

type INodeProvider interface {
	getBaseURL() string
}

type NodeInfo struct {
	URL string
}

type NodeProvider struct {
	SelectedNodeIndex int
	Nodes             []NodeInfo
}

func (nodeProvider NodeProvider) getBaseURL() string {
	return nodeProvider.Nodes[nodeProvider.SelectedNodeIndex].URL
}

// func checkNodeHealth(nodeProvider NodeProvider) {
//    Ping each node in nodeProvider.Nodes and set SelectedNodeIndex to the one with the best ping time.
// }
