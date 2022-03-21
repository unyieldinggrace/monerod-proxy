package endpoints

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/labstack/echo/v4"
	"gopkg.in/ini.v1"
)

func ConfigureMonerodProxyHandler(e *echo.Echo, cfg *ini.File) {
	e.GET(":monerodendpoint", func(c echo.Context) error {
		fmt.Println("Monerod Endpoint: " + c.Param("monerodendpoint"))
		baseURL := cfg.Section("").Key("node").Value()
		resp := forwardGETRequest(baseURL + c.Param("monerodendpoint"))
		return c.String(http.StatusOK, resp)
	})

	e.POST(":monerodendpoint", func(c echo.Context) error {
		requestBody, err := ioutil.ReadAll(c.Request().Body)

		reqDump := "POST Request received: " + c.Param("monerodendpoint") + "\n" + string(requestBody)
		fmt.Println(reqDump)

		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, reqDump)
		}

		baseURL := cfg.Section("").Key("node").Value()
		resp := forwardPOSTRequest(baseURL+c.Param("monerodendpoint"), requestBody)
		return c.String(http.StatusOK, resp)
	})
}

func forwardGETRequest(URL string) string {
	var content string
	err := requests.URL(URL).
		ContentType("application/json").
		ToString(&content).
		Fetch(context.Background())

	if err != nil {
		fmt.Println("Error calling backend node:")
		fmt.Println(err)
	}

	fmt.Println(content)
	return content
}

//
// net/http native
//
func forwardPOSTRequest(URL string, body []byte) string {
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	log.Printf("Response Body:\n%s", resBody)
	return string(resBody)
}

//
// Requests HTTP
//
// func forwardPOSTRequest(URL string, body []byte) string {
// 	// URL = "http://httpbin.org/post"
// 	// URL = "http://localhost:8081/"
// 	var content string
// 	fmt.Println(body)
// 	err := requests.URL(URL).
// 		// Client(&http.Client{
// 		// 	Transport: &http.Transport{
// 		// 		Proxy: http.ProxyFromEnvironment,
// 		// 		Dial: (&net.Dialer{
// 		// 			Timeout:   30 * time.Second,
// 		// 			KeepAlive: 30 * time.Second,
// 		// 		}).Dial,
// 		// 		TLSHandshakeTimeout: 10 * time.Second,
// 		// 		DisableCompression:  true,
// 		// 	}}).
// 		// Header("Blah-Blah", "Some Header Data").
// 		// Header("Content-Length", fmt.Sprint(cap(body))).
// 		// Header("Content-Length", "100").
// 		ContentType("application/json").
// 		BodyBytes([]byte("{\"method\":\"get_version\"}")).
// 		// BodyBytes(body).
// 		ToString(&content).
// 		Fetch(context.Background())

// 	if err != nil {
// 		fmt.Println("Error calling backend node:")
// 		fmt.Println(err)
// 	}

// 	fmt.Println(content)
// 	return content
// }

//
// Monaco HTTP
//
// func forwardPOSTRequest(URL string, body []byte) string {
// 	// URL = "http://httpbin.org/post"
// 	// URL = "http://localhost:8081/"
// 	var result interface{}
// 	c := request.Client{
// 		Context: context.Background(),
// 		URL:     URL,
// 		Method:  "POST",
// 		JSON:    body,
// 	}

// 	// c.PrintCURL()
// 	// ctx := c.initContext()
// 	// req := ctx.GetRequest()
// 	// cmd, err := curl.GetCommand(req)

// 	// if err != nil {
// 	// 	fmt.Println("Error generating curl command:")
// 	// 	fmt.Println(err)
// 	// }

// 	response := c.Send().Scan(&result)
// 	fmt.Println(response)

// 	if !response.OK() {
// 		fmt.Println("Error calling backend node:")
// 		fmt.Println(response.Error())
// 	}

// 	fmt.Println(response.Error())
// 	fmt.Println(response.String())
// 	return response.String()
// }

//
// Build CURL command
//
// func forwardPOSTRequest(URL string, body []byte) string {
// 	// URL = "http://httpbin.org/post"
// 	// URL = "http://localhost:8081/"
// 	cmd := generateCURLCommand(URL, body)
// 	// fmt.Println(cmd)
// 	// fmt.Println(cmd[0])
// 	// fmt.Println(cmd[1])
// 	// fmt.Println("Curl Command: ")
// 	// fmt.Printf("%v\n", cmd)

// 	// output, err := exec.Command(cmd[0], cmd[1:]...).Output()
// 	output, err := exec.Command("curl", "-d", "'{\"jsonrpc\":\"2.0\",\"id\":\"0\",\"method\":\"rpc_access_info\", \"params\":{\"client\":\"asdf\"}}'", "-H", "'Content-Type: application/json'", "'http://xmrnode.digitalcashtools.com:18081/json_rpc'").Output()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return string(output)
// }

// func generateCURLCommand(URL string, body []byte) []string {
// 	return []string{"curl", "-d", "'" + strings.Replace(string(body), "\n", "", -1) + "'", "-H", "'Content-Type: application/json'", "'" + URL + "'"}
// }
