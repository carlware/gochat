package cases

import (
	"context"
	"encoding/json"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/models"
	"github.com/google/uuid"
)

type OptsRoom struct {
	Room chatroom.Room
	BR   chatroom.BroadcastReceiver
}

func ListRoom(opts *OptsRoom) ([]*models.Room, error) {
	fetched, err := opts.Room.List(context.TODO())
	if err != nil {
		return []*models.Room{}, nil
	}

	return fetched, nil
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type RoomResponse struct {
	Type    string        `json:"type"`
	Message string        `json:"message"`
	Error   *models.Error `json:"error"`
}

func CreateRoom(opts *OptsRoom, req *CreateRoomRequest) (*models.Room, error) {
	room := &models.Room{
		Name: req.Name,
		ID:   uuid.New().String(),
	}
	cres := RoomResponse{
		Type:    "roomCreated",
		Message: req.Name,
	}
	encoded, _ := json.Marshal(cres)
	opts.BR.Broadcast(encoded)

	return opts.Room.Add(context.TODO(), room)
}
