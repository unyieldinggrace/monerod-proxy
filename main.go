package main

import (
	"context"
	"digitalcashtools/monerod-proxy/endpoints"
	"fmt"
	"os"

	"github.com/carlmjohnson/requests"
	"github.com/labstack/echo/v4"
	"gopkg.in/ini.v1"
)

func main() {
	// loadURL()
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to read config.ini")
		os.Exit(1)
	}

	http_port := cfg.Section("").Key("http_port").Value()
	fmt.Println("Port from config: ", http_port)

	e := echo.New()
	endpoints.ConfigurePing(e)
	endpoints.ConfigureMonerodProxyHandler(e)

	fmt.Println("Server running, test by visiting localhost:", http_port, "/ping")
	e.Logger.Fatal(e.Start(":" + http_port))
}

func loadURL() {
	var content string
	err := requests.URL("https://www.digitalcashtools.com").
		ToString(&content).
		Fetch(context.Background())

	if err != nil {
		fmt.Println("Error encountered")
	}

	fmt.Println(content)
}
