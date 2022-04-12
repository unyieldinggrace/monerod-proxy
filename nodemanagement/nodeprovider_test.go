package nodemanagement

import (
	"errors"
	"strings"
	"testing"
)

func TestAddNode(t *testing.T) {
	nodeProvider := getBasicNodeProvider()
	nodeProvider.AddNode("http://mynode.com:18081/")
	nodeProvider.AddNode("http://someothernode.com:18081/")

	numNodes := len(nodeProvider.Nodes)
	if numNodes != 2 {
		t.Errorf("Number of added nodes should be 2, got %d", numNodes)
	}
}

func TestGetBaseURL(t *testing.T) {
	nodeProvider := getBasicNodeProvider()
	expectedURL := "http://mynode.com:18081/"
	nodeProvider.AddNode(expectedURL)

	URL := nodeProvider.GetBaseURL()
	if URL != expectedURL {
		t.Errorf("Returned BaseURL: %s, should be %s", URL, expectedURL)
	}
}

func TestGetBaseURLWithTwoNodes(t *testing.T) {
	nodeProvider := getBasicNodeProvider()
	expectedURL := "http://mynode.com:18081/"
	nodeProvider.AddNode("http://someothernode.com:18081/")
	nodeProvider.AddNode(expectedURL)
	nodeProvider.SelectedNodeIndex = 1

	URL := nodeProvider.GetBaseURL()
	if URL != expectedURL {
		t.Errorf("Returned BaseURL: %s, should be %s", URL, expectedURL)
	}
}

func TestWhenCheckNodeHealthIsCalledAndTestRequestReturnsErrorThenAnyNodesAvailableIsFalse(t *testing.T) {
	httpRequestFunc := func(URL string) (string, int, error) {
		return "test", 500, errors.New("test error")
	}

	nodeProvider := getNodeProviderWithHTTPFunc(httpRequestFunc)
	nodeProvider.AddNode("http://mynode.com:18081/")
	nodeProvider.Nodes[0].PassedLastCheck = false

	nodeProvider.CheckNodeHealth()
	if nodeProvider.GetAnyNodesAvailable() {
		t.Errorf("GetAnyNodesAvailable() returned true, should be false")
	}
}

func TestWhenCheckNodeHealthIsCalledAndTestRequestReturnsSuccessThenAnyNodesAvailableIsTrue(t *testing.T) {
	httpRequestFunc := func(URL string) (string, int, error) {
		return "test", 200, nil
	}

	nodeProvider := getNodeProviderWithHTTPFunc(httpRequestFunc)
	nodeProvider.AddNode("http://mynode.com:18081/")
	nodeProvider.Nodes[0].PassedLastCheck = false

	nodeProvider.CheckNodeHealth()
	if !nodeProvider.GetAnyNodesAvailable() {
		t.Errorf("GetAnyNodesAvailable() returned false, should be true")
	}
}

func TestWhenCheckNodeHealthIsCalledAndCurrentNodeIsFailingThenSelectedNodeIndexIsShiftedToWorkingNode(t *testing.T) {
	failingNode := "http://mynode.com:18081/"
	workingNode := "http://someothernode.com:18081/"

	httpRequestFunc := func(URL string) (string, int, error) {
		if strings.Contains(URL, workingNode) {
			return "test", 200, nil
		} else {
			return "test", 500, errors.New("test error")
		}
	}

	nodeProvider := getNodeProviderWithHTTPFunc(httpRequestFunc)
	nodeProvider.AddNode(failingNode)
	nodeProvider.Nodes[0].PassedLastCheck = false
	nodeProvider.AddNode(workingNode)
	nodeProvider.Nodes[1].PassedLastCheck = false // will be changed by health check

	nodeProvider.CheckNodeHealth()
	result := nodeProvider.GetBaseURL()
	if result != workingNode {
		t.Errorf("Selected node should be %s, instead got %s", workingNode, result)
	}
}

func TestWhenCheckNodeHealthIsCalledAndFailingNodeNowSucceedsThenSelectedNodeIndexIsUnchanged(t *testing.T) {
	restoredNode := "http://mynode.com:18081/"
	workingNode := "http://someothernode.com:18081/"

	httpRequestFunc := func(URL string) (string, int, error) {
		return "test", 200, nil
	}

	nodeProvider := getNodeProviderWithHTTPFunc(httpRequestFunc)
	nodeProvider.AddNode(restoredNode)
	nodeProvider.Nodes[0].PassedLastCheck = false // will be changed by health check
	nodeProvider.AddNode(workingNode)
	nodeProvider.Nodes[1].PassedLastCheck = false // will be changed by health check

	nodeProvider.CheckNodeHealth()
	result := nodeProvider.GetBaseURL()
	if result != restoredNode {
		t.Errorf("Selected node should be %s, instead got %s", restoredNode, result)
	}
}

func TestWhenNodeFailureIsReportedThenNodeHealthGetsChecked(t *testing.T) {
	failingNode := "http://mynode.com:18081/"
	nodeHealthChecked := false

	httpRequestFunc := func(URL string) (string, int, error) {
		nodeHealthChecked = true
		return "test", 200, nil
	}

	nodeProvider := getNodeProviderWithHTTPFunc(httpRequestFunc)
	nodeProvider.AddNode(failingNode)
	nodeProvider.ReportNodeConnectionFailure()

	if !nodeHealthChecked {
		t.Errorf("Expected HTTP requests to check node health, but none occurred.")
	}
}

func getBasicNodeProvider() *NodeProvider {
	return getNodeProviderWithHTTPFunc(nil)
}

func getNodeProviderWithHTTPFunc(executeGETRequestFunc ExecuteGETRequestFunc) *NodeProvider {
	nodeProvider := &NodeProvider{
		SelectedNodeIndex:     0,
		Nodes:                 []NodeInfo{},
		AnyNodesAvailable:     true,
		executeGETRequestFunc: executeGETRequestFunc,
	}

	return nodeProvider
}
