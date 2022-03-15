package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
	"github.com/labstack/echo/v4"
)

type MonerodProxyHandlerResponse struct {
	MonerodEndpoint string `json:"monerodendpoint"`
}

func ConfigureMonerodProxyHandler(e *echo.Echo) {
	e.GET(":monerodendpoint", func(c echo.Context) error {
		fmt.Println("Monerod Endpoint: " + c.Param("monerodendpoint"))
		// resp := &MonerodProxyHandlerResponse{MonerodEndpoint: c.Param("monerodendpoint")}
		// return c.JSON(http.StatusOK, resp)
		resp := loadURL(c.Param("monerodendpoint"))
		return c.String(http.StatusOK, resp)
	})
}

func loadURL(endpoint string) string {
	var content string
	err := requests.URL("http://xmrnode.digitalcashtools.com:18081/" + endpoint).
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
