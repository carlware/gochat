package main

import (
	"context"
	"flag"

	"github.com/carlware/gochat/chatroom/cases"
	"github.com/carlware/gochat/common/auth"
	"github.com/carlware/gochat/dispatchers/rest"
	"github.com/carlware/gochat/dispatchers/websocket"
	"github.com/carlware/gochat/stockbot"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var addr = flag.String("addr", ":8080", "http server address")
var ctx = context.Background()

func main() {
	flag.Parse()
	e := echo.New()

	hub := websocket.NewHub()
	go hub.Run()
	go stockbot.RunRabbitMQ()
	go stockbot.Listen()
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
