package endpoints

import (
	"digitalcashtools/monerod-proxy/nodemanagement"
	"digitalcashtools/monerod-proxy/security"
	"encoding/json"
	"io/ioutil"
	"net/http"

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

type PasswordHolder struct {
	Password string `json:"Password"`
}

const adminPasswordRejectedMessage = "Admin password rejected."

func ConfigureAdminEndpoints(e *echo.Echo, passwordChecker security.IPasswordChecker, nodeProvider nodemanagement.INodeProvider) {
	e.GET("/proxy/api/status", func(c echo.Context) error {
		adminPasswordFound, err := checkAdminPassword(c, passwordChecker)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if !adminPasswordFound {
			return c.String(http.StatusForbidden, adminPasswordRejectedMessage)
		}

		return c.JSON(http.StatusOK, getStatusResponse(nodeProvider))
	})

	e.POST("/proxy/api/setnodeenabled", func(c echo.Context) error {
		adminPasswordFound, err := checkAdminPassword(c, passwordChecker)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if !adminPasswordFound {
			return c.String(http.StatusForbidden, adminPasswordRejectedMessage)
		}

		requestBody, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Debug(err)
			return c.String(http.StatusBadRequest, err.Error())
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

func checkAdminPassword(c echo.Context, passwordChecker security.IPasswordChecker) (bool, error) {
	requestBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return false, err
	}

	var passwordFromBody PasswordHolder
	err = json.Unmarshal(requestBody, &passwordFromBody)
	if err != nil {
		return false, err
	}

	return passwordChecker.CheckAdminPassword(passwordFromBody.Password), nil
}
