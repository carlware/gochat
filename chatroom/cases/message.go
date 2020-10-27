package cases

import (
	"context"
	"sort"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/models"
)

type OptsMessage struct {
	Message chatroom.Message
}

func ListMessages(opts *OptsMessage, rid string) ([]*models.Message, error) {
	fetched, err := opts.Message.List(context.TODO())
	if err != nil {
		return []*models.Message{}, nil
	}
	sort.Slice(fetched, func(i, j int) bool { return fetched[i].Created.Before(fetched[j].Created) })

	msgs := []*models.Message{}
	counter := 50
	for _, msg := range fetched {
		if msg.RID == rid {
			msgs = append(msgs, msg)
			counter--
		}
		if counter == 0 {
			break
		}
	}
	return msgs, nil
}
