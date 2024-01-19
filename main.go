package main

import (
	"btc/web"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"runtime"
)

var (
	version       = "main"
	gitCommitHash = ""
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:8000", "")
	flag.Parse()

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
