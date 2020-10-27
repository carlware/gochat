package cases

import (
	"context"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/models"
	"github.com/google/uuid"
)

type OptsRoom struct {
	Room chatroom.Room
}

func ListRoom(opts *OptsRoom) ([]*models.Room, error) {
	fetched, err := opts.Room.List(context.TODO())
	if err != nil {
		return []*models.Room{}, nil
	}

	// sort.Slice(fetched, func(i, j int) bool { return fetched[i].Created < fetched[j].Created })

	return fetched, nil
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

func CreateRoom(opts *OptsRoom, req *CreateRoomRequest) (*models.Room, error) {
	return opts.Room.Add(context.TODO(), &models.Room{
		Name: req.Name,
		ID:   uuid.New().String(),
	})
}
