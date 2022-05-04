package endpoints

import (
	"digitalcashtools/monerod-proxy/nodemanagement"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type StatusResponse struct {
	CurrentNode    string   `json:"CurrentNode"`
	AvailableNodes []string `json:"AvailableNodes"`
}

type SetNodeEnabledRequestBody struct {
	NodeURL string `json:"NodeURL"`
	Enabled bool   `json:"Enabled"`
}

func ConfigureAdminEndpoints(e *echo.Echo, nodeProvider nodemanagement.INodeProvider) {
	e.GET("/proxy/api/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, getStatusResponse(nodeProvider))
	})

	e.POST("/proxy/api/setnodeenabled", func(c echo.Context) error {
		requestBody, err := ioutil.ReadAll(c.Request().Body)
		reqDump := time.Now().Format(time.RFC3339) + " " + c.RealIP() + " POST Request received: disablenode"

		if err != nil {
			log.Debug(err)
			return c.String(http.StatusBadRequest, reqDump)
		}

		var requestStruct SetNodeEnabledRequestBody
		err = json.Unmarshal(requestBody, &requestStruct)
		if err != nil {
			return c.String(http.StatusBadRequest, "Could not parse request body.")
		}

		nodeFound := nodeProvider.SetNodeEnabled(requestStruct.NodeURL, requestStruct.Enabled)
		if !nodeFound {
			return c.String(http.StatusBadRequest, "No node found with URL "+requestStruct.NodeURL)
		}

		nodeProvider.CheckNodeHealth()
		return c.JSON(http.StatusOK, getStatusResponse(nodeProvider))
	})
}

func getStatusResponse(nodeProvider nodemanagement.INodeProvider) *StatusResponse {
	return &StatusResponse{
		CurrentNode:    nodeProvider.GetBaseURL(),
		AvailableNodes: nodeProvider.GetAvailableNodes(),
	}
}
