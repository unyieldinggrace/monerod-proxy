package endpoints

import (
	"digitalcashtools/monerod-proxy/httpclient"
	"digitalcashtools/monerod-proxy/nodemanagement"

	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func ConfigureMonerodProxyHandler(e *echo.Echo, nodeProvider nodemanagement.INodeProvider) {
	e.GET(":monerodendpoint", func(c echo.Context) error {
		if !nodeProvider.GetAnyNodesAvailable() {
			return getNoNodesAvailableResponse(c)
		}

		baseURL := nodeProvider.GetBaseURL()
		resp, httpStatus, err := forwardGETRequest(baseURL + c.Param("monerodendpoint"))
		if err != nil {
			nodeProvider.ReportNodeConnectionFailure()
		}

		log.Debug(time.Now().Format(time.RFC3339), " ", c.RealIP(), " GET Request Received: ", c.Param("monerodendpoint"), "\tResponse Code: ", httpStatus)
		return c.String(httpStatus, resp)
	})

	e.POST(":monerodendpoint", func(c echo.Context) error {
		if !nodeProvider.GetAnyNodesAvailable() {
			return getNoNodesAvailableResponse(c)
		}

		requestBody, err := ioutil.ReadAll(c.Request().Body)
		reqDump := time.Now().Format(time.RFC3339) + " " + c.RealIP() + " POST Request received: " + c.Param("monerodendpoint")

		if err != nil {
			log.Debug(err)
			return c.String(http.StatusBadRequest, reqDump)
		}

		baseURL := nodeProvider.GetBaseURL()
		resp, httpStatus, err := forwardPOSTRequest(baseURL+c.Param("monerodendpoint"), requestBody)
		if err != nil {
			nodeProvider.ReportNodeConnectionFailure()
		}

		log.Debug(reqDump, "\tResponse Code: ", httpStatus)
		return c.String(httpStatus, resp)
	})
}

func forwardGETRequest(URL string) (string, int, error) {
	return httpclient.ExecuteGETRequest(URL)
}

func forwardPOSTRequest(URL string, body []byte) (string, int, error) {
	return httpclient.ExecutePOSTRequest(URL, body)
}

func getNoNodesAvailableResponse(c echo.Context) error {
	return c.String(500, "No nodes available.")
}
