package main

import (
	"btc/web"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:8000", "")
	flag.Parse()

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.HideBanner = true

	// addresses
	app.Match([]string{"GET", "POST"}, "/address/new", web.NewAddress)
	app.Match([]string{"GET", "POST"}, "/address/parse", web.ParseAddress)
	app.Match([]string{"GET", "POST"}, "/address/new-multi-sig", web.NewMultiSigAddress)

	// transactions
	app.Match([]string{"GET", "POST"}, "/transaction/decode", web.DecodeTransaction)
	app.Match([]string{"GET", "POST"}, "/transaction/create", web.CreateRawTransaction)
	app.Match([]string{"GET", "POST"}, "/transaction/sign", web.SignTransaction)

	app.Logger.Fatal(app.Start(addr))
}
