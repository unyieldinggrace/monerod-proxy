package nodemanagement

import (
	"testing"
)

func TestGetBaseURK(t *testing.T) {
	nodeProvider := getBasicNodeProvider()
	expectedURL := "http://mynode.com:18081"
	nodeProvider.AddNode(expectedURL)

	URL := nodeProvider.GetBaseURL()
	if URL != expectedURL {
		t.Errorf("Returned BaseURL: %s, should be %s", URL, expectedURL)
	}
}

func getBasicNodeProvider() *NodeProvider {
	nodeProvider := &NodeProvider{
		SelectedNodeIndex:     0,
		Nodes:                 []NodeInfo{},
		AnyNodesAvailable:     true,
		executeGETRequestFunc: nil,
	}

	return nodeProvider
}
