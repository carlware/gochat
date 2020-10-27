package memorydb

import (
	"context"
	"errors"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/models"
)

type message struct {
}

func NewMessage() chatroom.Message {
	return &message{}
}

var DBMessage = map[string]*models.Message{}

func (p *message) Add(ctx context.Context, message *models.Message) (*models.Message, error) {
	DBMessage[message.ID] = message
	return message, nil
}

func (p *message) Remove(ctx context.Context, message *models.Message) (*models.Message, error) {
	if _, ok := DBMessage[message.ID]; ok {
		delete(DBMessage, message.ID)
		return message, nil
	} else {
		return nil, errors.New("A resource with this ID does not exists")
	}
}

func (p *message) Get(ctx context.Context, messageID string) (*models.Message, error) {
	if val, ok := DBMessage[messageID]; ok {
		return val, nil
	} else {
		return nil, errors.New("A resource with this ID does not exists")
	}
}

func (p *message) Update(ctx context.Context, message *models.Message) (*models.Message, error) {
	if _, ok := DBMessage[message.ID]; ok {
		DBMessage[message.ID] = message
		return message, nil
	} else {
		return nil, errors.New("A resource with this ID does not exists")
	}
}

func (p *message) List(ctx context.Context) ([]*models.Message, error) {
	messages := []*models.Message{}
	for _, message := range DBMessage {
		messages = append(messages, message)
	}
	return messages, nil
}
