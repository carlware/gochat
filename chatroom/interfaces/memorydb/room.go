package memorydb

import (
	"context"
	"errors"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/models"
)

// Warining: this does not support currency, just for demo

type room struct {
}

func NewRoom() chatroom.Room {
	return &room{}
}

var DBRoom = map[string]*models.Room{}

func (p *room) Add(ctx context.Context, room *models.Room) (*models.Room, error) {
	DBRoom[room.ID] = room
	return room, nil
}

func (p *room) Remove(ctx context.Context, room *models.Room) (*models.Room, error) {
	if _, ok := DBRoom[room.ID]; ok {
		delete(DBRoom, room.ID)
		return room, nil
	} else {
		return nil, errors.New("A resource with this ID does not exists")
	}
}

func (p *room) Get(ctx context.Context, roomID string) (*models.Room, error) {
	if val, ok := DBRoom[roomID]; ok {
		return val, nil
	} else {
		return nil, errors.New("A resource with this ID does not exists")
	}
}

func (p *room) Update(ctx context.Context, room *models.Room) (*models.Room, error) {
	if _, ok := DBRoom[room.ID]; ok {
		DBRoom[room.ID] = room
		return room, nil
	} else {
		return nil, errors.New("A resource with this ID does not exists")
	}
}

func (p *room) List(ctx context.Context) ([]*models.Room, error) {
	rooms := []*models.Room{}
	for _, room := range DBRoom {
		rooms = append(rooms, room)
	}
	return rooms, nil
}
