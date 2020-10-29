package rest

import (
	"errors"
	"net/http"

	"github.com/carlware/gochat/chatroom/cases"
	"github.com/carlware/gochat/chatroom/models"
	"github.com/labstack/echo/v4"
)

func (ctr *controller) ListMessage(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return errors.New(`URL should contain Chatroom ID`)
	}

	rooms, err := cases.ListMessages(&cases.OptsMessage{
		Message: ctr.messagedb,
	}, idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Code:    "server",
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, rooms)
}
