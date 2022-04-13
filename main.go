package main

import (
	"digitalcashtools/monerod-proxy/endpoints"
	"digitalcashtools/monerod-proxy/nodemanagement"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Failed to read config.ini")
	}

	logLevel := cfg.Section("").Key("log_level").Value()
	switch logLevel {
	case "Trace":
		log.SetLevel(log.TraceLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.Info("Log level from config not recognised, defaulting to Info level.")
	}

	httpPort := cfg.Section("").Key("http_port").Value()

	e := echo.New()
	endpoints.ConfigurePing(e)

	nodeProvider := nodemanagement.LoadNodeProviderFromConfig(cfg)
	endpoints.ConfigureMonerodProxyHandler(e, nodeProvider)

	setUpNodeHealthCheckTicker(cfg, nodeProvider)

	nodeProvider.CheckNodeHealth()
	log.Info("Selected node: ", nodeProvider.GetBaseURL())

	e.GET("*", func(c echo.Context) error {
		reqDump := time.Now().Format(time.RFC3339) + " GET Request received: " + c.Path() + c.QueryString()
		log.Debug(reqDump)
		return c.String(http.StatusOK, reqDump)
	})

	e.POST("*", func(c echo.Context) error {
		reqDump := time.Now().Format(time.RFC3339) + " POST Request received: " + c.Path() + c.QueryString()
		log.Debug(reqDump)
		return c.String(http.StatusOK, reqDump)
	})

	log.Info("Server running, test by visiting localhost:", httpPort, "/ping")

	// TODO: Add TLS support
	e.Logger.Fatal(e.Start(":" + httpPort))
}

func setUpNodeHealthCheckTicker(cfg *ini.File, nodeProvider nodemanagement.INodeProvider) {
	secondsBetweenHealthChecks, err := cfg.Section("").Key("seconds_between_health_checks").Int()
	if err != nil {
		secondsBetweenHealthChecks = 10 * 60
	}

	log.Info(fmt.Sprintf("Performing node health check every %d seconds.", secondsBetweenHealthChecks))

	healthCheckTicker := time.NewTicker(time.Duration(secondsBetweenHealthChecks) * time.Second)

	go func() {
		for {
			<-healthCheckTicker.C
			log.Debug("Checking node health...")

			previousNode := nodeProvider.GetBaseURL()

			nodeProvider.CheckNodeHealth()

			newSelectedNode := nodeProvider.GetBaseURL()

			if newSelectedNode != previousNode {
				log.Info("Switched to new node:", newSelectedNode)
			}
		}
	}()
}
