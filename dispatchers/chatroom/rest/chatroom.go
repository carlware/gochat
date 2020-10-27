package rest

import (
	"net/http"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/cases"
	"github.com/carlware/gochat/chatroom/models"
	"github.com/labstack/echo/v4"
)

type controller struct {
	roomdb    chatroom.Room
	messagedb chatroom.Message
}

func NewChatController(roomdb chatroom.Room, messagedb chatroom.Message) *controller {
	return &controller{
		roomdb:    roomdb,
		messagedb: messagedb,
	}
}

func (ctr *controller) ListRoom(c echo.Context) error {
	rooms, err := cases.ListRoom(&cases.OptsRoom{
		Room: ctr.roomdb,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Code:    "Server",
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, rooms)
}

func (ctr *controller) CreateRoom(c echo.Context) error {
	req := cases.CreateRoomRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Code:    "login",
			Message: err.Error(),
		})
	}

	room, err := cases.CreateRoom(&cases.OptsRoom{
		Room: ctr.roomdb,
	}, &req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Code:    "Server",
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, room)
}
