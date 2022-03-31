package main

import (
	"digitalcashtools/monerod-proxy/endpoints"
	"digitalcashtools/monerod-proxy/nodemanagement"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to read config.ini")
		os.Exit(1)
	}

	http_port := cfg.Section("").Key("http_port").Value()
	fmt.Println("Port from config: ", http_port)

	e := echo.New()
	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	fmt.Println("Request Body Dump")
	// 	fmt.Println(string(reqBody))
	// }))
	endpoints.ConfigurePing(e)
	// endpoints.ConfigureMonerodProxyHandler(e, cfg)

	// Create NodeProvider instance
	// Load nodes from config
	// Start timer to periodically run node health checks
	//
	endpoints.ConfigureMonerodProxyHandler(e, nodemanagement.LoadNodeProviderFromConfig(cfg))

	e.GET("*", func(c echo.Context) error {
		reqDump := time.Now().Format(time.RFC3339) + " GET Request received: " + c.Path() + c.QueryString()
		fmt.Println(reqDump)
		return c.String(http.StatusOK, reqDump)
	})

	e.POST("*", func(c echo.Context) error {
		reqDump := time.Now().Format(time.RFC3339) + " POST Request received: " + c.Path() + c.QueryString()
		fmt.Println(reqDump)
		return c.String(http.StatusOK, reqDump)
	})

	fmt.Println("Server running, test by visiting localhost:", http_port, "/ping")
	e.Logger.Fatal(e.Start(":" + http_port))
}
