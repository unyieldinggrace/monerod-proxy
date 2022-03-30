package endpoints

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/ini.v1"
)

func ConfigureMonerodProxyHandler(e *echo.Echo, cfg *ini.File) {
	e.GET(":monerodendpoint", func(c echo.Context) error {
		fmt.Print(time.Now().Format(time.RFC3339) + "\tGET Request Received: " + c.Param("monerodendpoint"))
		baseURL := cfg.Section("").Key("node").Value()
		resp, httpStatus := forwardGETRequest(baseURL + c.Param("monerodendpoint"))
		fmt.Println("\tResponse Code: ", httpStatus)
		return c.String(httpStatus, resp)
	})

	e.POST(":monerodendpoint", func(c echo.Context) error {
		requestBody, err := ioutil.ReadAll(c.Request().Body)

		//reqDump := "POST Request received: " + c.Param("monerodendpoint") + "\n" + string(requestBody)
		reqDump := time.Now().Format(time.RFC3339) + "\tPOST Request received: " + c.Param("monerodendpoint")
		fmt.Print(reqDump)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, reqDump)
		}

		baseURL := cfg.Section("").Key("node").Value()
		resp, httpStatus := forwardPOSTRequest(baseURL+c.Param("monerodendpoint"), requestBody)
		fmt.Println("\tResponse Code: ", httpStatus)
		return c.String(httpStatus, resp)
	})
}

func forwardGETRequest(URL string) (string, int) {
	req, _ := http.NewRequest("GET", URL, nil)
	return executeRequest(req)
}

func forwardPOSTRequest(URL string, body []byte) (string, int) {
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	return executeRequest(req)
}

func executeRequest(req *http.Request) (string, int) {
	res, _ := http.DefaultClient.Do(req)
	if !getStatusCodeSuccessful(res) {
		return fmt.Sprint(res.StatusCode, ": ", http.StatusText(res.StatusCode)), res.StatusCode
	}
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	// log.Printf("Response Body:\n%s", resBody)
	return string(resBody), res.StatusCode
}

func getStatusCodeSuccessful(res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode <= 299
}
