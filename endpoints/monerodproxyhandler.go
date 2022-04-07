package endpoints

import (
	"digitalcashtools/monerod-proxy/httpclient"
	"digitalcashtools/monerod-proxy/nodemanagement"

	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

		fmt.Println(time.Now().Format(time.RFC3339)+" GET Request Received: "+c.Param("monerodendpoint")+"\tResponse Code: ", httpStatus)
		return c.String(httpStatus, resp)
	})

	e.POST(":monerodendpoint", func(c echo.Context) error {
		if !nodeProvider.GetAnyNodesAvailable() {
			return getNoNodesAvailableResponse(c)
		}

		requestBody, err := ioutil.ReadAll(c.Request().Body)

		//reqDump := "POST Request received: " + c.Param("monerodendpoint") + "\n" + string(requestBody)
		reqDump := time.Now().Format(time.RFC3339) + " POST Request received: " + c.Param("monerodendpoint")

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, reqDump)
		}

		baseURL := nodeProvider.GetBaseURL()
		resp, httpStatus, err := forwardPOSTRequest(baseURL+c.Param("monerodendpoint"), requestBody)
		if err != nil {
			nodeProvider.ReportNodeConnectionFailure()
		}

		fmt.Println(reqDump, "\tResponse Code: ", httpStatus)
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
