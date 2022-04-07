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

	nodeProvider := nodemanagement.LoadNodeProviderFromConfig(cfg)
	endpoints.ConfigureMonerodProxyHandler(e, nodeProvider)

	setUpNodeHealthCheckTicker(cfg, nodeProvider)

	nodeProvider.CheckNodeHealth()
	fmt.Println("Selected node: ", nodeProvider.GetBaseURL())

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

func setUpNodeHealthCheckTicker(cfg *ini.File, nodeProvider nodemanagement.INodeProvider) {
	secondsBetweenHealthChecks, err := cfg.Section("").Key("seconds_between_health_checks").Int()
	if err != nil {
		secondsBetweenHealthChecks = 10 * 60
	}

	fmt.Println("Performing node health check every ", secondsBetweenHealthChecks, " seconds")

	healthCheckTicker := time.NewTicker(time.Duration(secondsBetweenHealthChecks) * time.Second)

	go func() {
		for {
			<-healthCheckTicker.C
			fmt.Println("Checking node health...")
			nodeProvider.CheckNodeHealth()
		}
	}()
}
