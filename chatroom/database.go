package chatroom

import (
	"context"

	"github.com/carlware/gochat/chatroom/models"
)

// Room
type Room interface {
	Add(context.Context, *models.Room) (*models.Room, error)
	List(context.Context) ([]*models.Room, error)
}

// Message
type Message interface {
	Add(context.Context, *models.Message) (*models.Message, error)
	List(context.Context) ([]*models.Message, error)
}
