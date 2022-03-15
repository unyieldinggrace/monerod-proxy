package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MonerodProxyHandlerResponse struct {
	MonerodEndpoint string `json:"monerodendpoint"`
}

func ConfigureMonerodProxyHandler(e *echo.Echo) {
	e.GET(":monerodendpoint", func(c echo.Context) error {
		resp := &MonerodProxyHandlerResponse{MonerodEndpoint: c.Param("monerodendpoint")}
		return c.JSON(http.StatusOK, resp)
	})
}
