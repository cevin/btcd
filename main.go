package main

import (
	"btc/config"
	"btc/web"
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	version       = "main"
	gitCommitHash = "latest"
)

func main() {
	var addr string
	var v bool
	var tNet string
	flag.StringVar(&addr, "addr", "localhost:8000", "")
	flag.StringVar(&tNet, "net", "mainnet", "")
	flag.BoolVar(&v, "version", false, "Display version info")
	flag.Parse()

	if v {
		fmt.Printf(
			"Version: %s \n"+
				"Git Commit Hash: %s \n",
			version,
			gitCommitHash,
		)
		return
	}

	switch tNet {
	case "mainnet":
		config.NET = &chaincfg.MainNetParams
	case "testnet":
		config.NET = &chaincfg.TestNet3Params
	default:
		log.Fatal("Unknown network type: ", tNet)
	}

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.HideBanner = true

	app.RouteNotFound("/*", web.DefaultNotFound)

	// info
	app.Any("/", func(c echo.Context) error {
		return c.JSON(200, struct {
			Code      int    `json:"code"`
			Version   string `json:"version"`
			GoVersion string `json:"go-version"`
			GitHash   string `json:"git-commit-hash"`
		}{
			Code:      200,
			Version:   version,
			GoVersion: runtime.Version(),
			GitHash:   gitCommitHash,
		})
	})

	// addresses
	app.Match([]string{"GET", "POST"}, "/address/new", web.NewAddress)
	app.Match([]string{"GET", "POST"}, "/address/parse", web.ParseAddress)
	app.Match([]string{"GET", "POST"}, "/address/new-multi-sig", web.NewMultiSigAddress)

	// transactions
	app.Match([]string{"GET", "POST"}, "/transaction/decode", web.DecodeTransaction)
	app.POST("/transaction/create", web.CreateRawTransaction)
	app.POST("/transaction/sign", web.SignTransaction)
	app.POST("/transaction/create-and-sign", web.CreateAndSignTransaction)

	app.Logger.Fatal(app.Start(addr))
}
