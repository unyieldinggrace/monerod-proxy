package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingResponse struct {
	Message string `json:"message"`
}

func ConfigurePing(e *echo.Echo) {
	e.GET("/proxy/ping", func(c echo.Context) error {
		pingResponse := &PingResponse{Message: "pong"}
		return c.JSON(http.StatusOK, pingResponse)
	})
}
