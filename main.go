package main

import (
	"context"
	"flag"

	"github.com/carlware/gochat/common/config"
	"github.com/carlware/gochat/dispatchers/accounts"
	"github.com/carlware/gochat/dispatchers/chatroom"
	"github.com/carlware/gochat/stockbot"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var addr = flag.String("addr", ":8080", "http server address")
var cfgFile = flag.String("file", "", "conf file")
var ctx = context.Background()

func main() {
	flag.Parse()

	// configure echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{"POST", "GET"},
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// dispatch react app
	e.Static("/", "./web/build")

	// get default configuration or use file configuration
	cfg := &config.Configuration{}
	config.Load(cfg, "GOCHAT", *cfgFile)

	// run micro-service chat
	chatroom.RunMicroChatroom(e, cfg)

	// run micro-service accounts
	accounts.RunMicroAccounts(e)

	// stockbot
	go stockbot.Run(cfg)

	go e.Logger.Fatal(e.Start(*addr))
}
