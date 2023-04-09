package chatroom

import (
	"context"

	"github.com/carlware/gochat/chatroom/cases"
	"github.com/carlware/gochat/chatroom/cmds"
	"github.com/carlware/gochat/chatroom/interfaces/memorydb"
	"github.com/carlware/gochat/chatroom/models"
	"github.com/carlware/gochat/common/auth"
	"github.com/carlware/gochat/common/config"
	"github.com/carlware/gochat/common/mq/interfaces/rabbitmq"
	"github.com/carlware/gochat/dispatchers/chatroom/rest"
	"github.com/carlware/gochat/dispatchers/websocket"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RunMicroChatroom wire up everything needed
func RunMicroChatroom(e *echo.Echo, cfg *config.Configuration) {
	dbRoom := memorydb.NewRoom()
	dbMessage := memorydb.NewMessage()

	// adding first room
	dbRoom.Add(context.TODO(), &models.Room{
		ID:   uuid.New().String(),
		Name: "general",
	})

	// run websocket listener
	hub := websocket.NewHub()
	go hub.Run()

	ctrl := rest.NewChatController(dbRoom, dbMessage, hub)

	// rooms API
	e.GET("/rooms", ctrl.ListRoom, auth.IsLoggedIn)
	e.POST("/rooms", ctrl.CreateRoom, auth.IsLoggedIn)

	// messages API
	e.GET("/messages/:id", ctrl.ListMessage, auth.IsLoggedIn)

	// Websocket API
	e.GET("/ws", ctrl.WS)

	// Create rabbit queue
	queue, err := rabbitmq.NewServer(cfg.RabbiMQ.Host)
	if err != nil {
		panic(err)
	}

	// start commands processor
	cProcesor := cmds.NewCommandProcessor(queue)
	go cProcesor.Run()

	// listen incoming messages from websocket and process commands from it
	mListener := cases.NewMessageListener(hub, cProcesor, dbMessage)
	mListener.Listen()
}
