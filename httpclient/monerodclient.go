package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ExecuteGETRequest(URL string) (string, int, error) {
	req, _ := http.NewRequest("GET", URL, nil)
	return executeRequest(req)
}

func ExecutePOSTRequest(URL string, body []byte) (string, int, error) {
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	return executeRequest(req)
}

func executeRequest(req *http.Request) (string, int, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return "Monerod-proxy encountered an error connecting to a backend node. Please re-try your request and monerod-proxy will attempt to use an alternate node for fulfillment.\n", 500, err
	}

	if !getStatusCodeSuccessful(res) {
		return fmt.Sprint(res.StatusCode, ": ", http.StatusText(res.StatusCode)), res.StatusCode, nil
	}

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	return string(resBody), res.StatusCode, nil
}

func getStatusCodeSuccessful(res *http.Response) bool {
	if res == nil {
		return false
	}

	return res.StatusCode >= 200 && res.StatusCode <= 299
}
