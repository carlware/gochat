package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/carlware/gochat/chatroom/cases"
	"github.com/carlware/gochat/chatroom/interfaces/rabbitmq"
	"github.com/carlware/gochat/common/auth"
	"github.com/carlware/gochat/common/config"
	"github.com/carlware/gochat/dispatchers/rest"
	"github.com/carlware/gochat/dispatchers/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var addr = flag.String("addr", ":8080", "http server address")
var cfgFile = flag.String("file", "", "conf file")
var ctx = context.Background()

func main() {
	flag.Parse()
	e := echo.New()

	cfg := &config.Configuration{}
	config.Load(cfg, "GOCHAT", *cfgFile)
	fmt.Println(cfg)

	hub := websocket.NewHub()
	mQ, err := rabbitmq.NewServer(cfg.RabbiMQ.Host, cfg.RabbiMQ.CommandReqQueue)
	if err != nil {
		panic(err)
	}

	go hub.Run()
	go mQ.Listen()
	go cases.ListenMessages(&cases.Options{
		BroadcastReceiver: hub,
	})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./web/build")
	e.GET("/ws", func(c echo.Context) error {
		websocket.ServeWs(hub, c.Response(), c.Request())
		return nil
	})
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{"POST", "GET"},
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/profiles", rest.ListProfiles, auth.IsLoggedIn)
	e.POST("/login", rest.Login)

	go e.Logger.Fatal(e.Start(*addr))
}
