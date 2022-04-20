package main

import (
	"fmt"
	"net/http"

	"github.com/cheqd/cheqd-did-resolver/services"
	"github.com/cheqd/cheqd-did-resolver/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"strings"
)

func main() {
	viper.SetConfigFile("config.yaml")
	err1 := viper.ReadInConfig()

	viper.SetConfigFile(".env")
	err2 := viper.ReadInConfig()
	if err1 != nil && err2 != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n Fatal error config file: %s\n", err1.Error(), err2.Error()))
	}

	viper.AutomaticEnv()

	didResolutionPath := viper.GetString("resolverPath")
	didResolutionListener := viper.GetString("listener")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//setup
	networks := viper.GetStringMapString("networks")
	ledgerService := services.NewLedgerService()
	for network, url := range networks {
		e.StdLogger.Println(network)
		ledgerService.RegisterLedger(network, url)
	}
	requestService := services.NewRequestService(ledgerService)

	// Routes
	e.GET(didResolutionPath, func(c echo.Context) error {
		didUrl := c.Request().URL.String()
		fmt.Println(c.Request().URL.String())
		accept := strings.Split(c.Request().Header.Get("accept"), ";")[0]
		var acceptOption types.ContentType
		if strings.Contains(accept, string(types.JSONLD)) {
			acceptOption = types.DIDJSONLD
		} else {
			acceptOption = types.DIDJSON
		}
		e.StdLogger.Println("get did")
		responseBody, err := requestService.ProcessDIDRequest(didUrl, types.ResolutionOption{Accept: acceptOption})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Response().Header().Set(string(acceptOption), accept)
		return c.JSONBlob(http.StatusOK, []byte(responseBody))
	})

	e.Logger.Fatal(e.Start(didResolutionListener))
}
