package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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
		fmt.Println(err)
		return "Error connecting to backend node.\n", 500, err
	}

	if !getStatusCodeSuccessful(res) {
		return fmt.Sprint(res.StatusCode, ": ", http.StatusText(res.StatusCode)), res.StatusCode, nil
	}

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	// log.Printf("Response Body:\n%s", resBody)
	return string(resBody), res.StatusCode, nil
}

func getStatusCodeSuccessful(res *http.Response) bool {
	if res == nil {
		return false
	}

	return res.StatusCode >= 200 && res.StatusCode <= 299
}
