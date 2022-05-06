package endpoints

import (
	"digitalcashtools/monerod-proxy/nodemanagement"
	"digitalcashtools/monerod-proxy/security"
	"encoding/json"
	"fmt"
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

type InnerRequestHandler func(echo.Context, []byte) error

const adminPasswordRejectedMessage = "Admin password rejected."

func ConfigureAdminEndpoints(e *echo.Echo, passwordChecker security.IPasswordChecker, nodeProvider nodemanagement.INodeProvider) {
	e.POST("/proxy/api/status", func(c echo.Context) error {
		return handleRequestWrapper(c, passwordChecker, func(c echo.Context, requestBody []byte) error {
			return c.JSON(http.StatusOK, getStatusResponse(nodeProvider))
		})
	})

	e.POST("/proxy/api/setnodeenabled", func(c echo.Context) error {
		return handleRequestWrapper(c, passwordChecker, func(c echo.Context, requestBody []byte) error {
			var requestStruct SetNodeEnabledRequestBody
			err := json.Unmarshal(requestBody, &requestStruct)
			if err != nil {
				errorMessage := fmt.Sprintf("Could not parse request body.\n%s\n%s\n", requestBody, err.Error())
				return c.String(http.StatusBadRequest, errorMessage)
			}

			nodeFound := nodeProvider.SetNodeEnabled(requestStruct.NodeURL, requestStruct.Enabled)
			if !nodeFound {
				return c.String(http.StatusBadRequest, "No node found with URL "+requestStruct.NodeURL)
			}

			nodeProvider.CheckNodeHealth()
			return c.JSON(http.StatusOK, getStatusResponse(nodeProvider))
		})
	})

	e.POST("/proxy/api/generatepasswordhash", func(c echo.Context) error {
		return handleRequestWrapper(c, nil, func(c echo.Context, requestBody []byte) error {
			var requestStruct PasswordHolder
			err := json.Unmarshal(requestBody, &requestStruct)
			if err != nil {
				return c.String(http.StatusBadRequest, "Could not parse request body.")
			}

			result, err := passwordChecker.GeneratePasswordHash(requestStruct.Password)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			return c.String(http.StatusOK, result)
		})
	})
}

func handleRequestWrapper(c echo.Context, passwordChecker security.IPasswordChecker, innerRequestHandler InnerRequestHandler) error {
	requestBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Debug(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	if passwordChecker != nil {
		adminPasswordFound, err := checkAdminPassword(requestBody, passwordChecker)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if !adminPasswordFound {
			return c.String(http.StatusForbidden, adminPasswordRejectedMessage)
		}
	}

	return innerRequestHandler(c, requestBody)
}

func getStatusResponse(nodeProvider nodemanagement.INodeProvider) *StatusResponse {
	return &StatusResponse{
		CurrentNode:    nodeProvider.GetBaseURL(),
		AvailableNodes: nodeProvider.GetAvailableNodes(),
	}
}

func checkAdminPassword(requestBody []byte, passwordChecker security.IPasswordChecker) (bool, error) {
	var passwordFromBody PasswordHolder
	err := json.Unmarshal(requestBody, &passwordFromBody)
	if err != nil {
		return false, err
	}

	return passwordChecker.CheckAdminPassword(passwordFromBody.Password), nil
}
