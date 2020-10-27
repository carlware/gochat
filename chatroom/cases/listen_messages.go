package cases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/cmds"
	"github.com/carlware/gochat/chatroom/models"
	"github.com/google/uuid"
)

type mServer struct {
	br        chatroom.BroadcastReceiver
	cp        cmds.CommandProcesor
	messagedb chatroom.Message
}

type MessageRequest struct {
	RID     string `json:"rid"`
	UID     string `json:"uid"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type MessageResponse struct {
	RID     string        `json:"rid"`
	UID     string        `json:"uid"`
	Type    string        `json:"type"`
	Message string        `json:"message"`
	Error   *models.Error `json:"error"`
}

type CommandExtra struct {
	RID string `json:"rid"`
	UID string `json:"uid"`
}

func NewMessageListener(br chatroom.BroadcastReceiver, cp cmds.CommandProcesor, mdb chatroom.Message) *mServer {
	return &mServer{
		br:        br,
		cp:        cp,
		messagedb: mdb,
	}
}

func (s *mServer) Listen() {
	// listen
	messages, _ := s.br.Receive()

	// Listen incoming messages
	go func() {
		for msg := range messages {
			messageProccesor(msg, s.br, s.cp, s.messagedb)
		}
	}()

	// Listen incoming responses from command procesor
	go func() {
		results := s.cp.Results()
		for result := range results {
			fmt.Println("message command response procesor", result)
			extra := CommandExtra{}
			_ = json.Unmarshal([]byte(result.Extra), &extra)
			response := &MessageResponse{
				Message: result.Result,
				Type:    "command",
				RID:     extra.RID,
				UID:     extra.UID,
				Error:   result.Error,
			}
			encoded, _ := json.Marshal(response)
			s.br.Broadcast(encoded)
			addMessage(s.messagedb, extra.RID, result.Result, "bot")
		}
	}()
}

func messageProccesor(raw []byte, br chatroom.BroadcastReceiver, cp cmds.CommandProcesor, mdb chatroom.Message) {
	req := MessageRequest{}

	err := json.Unmarshal(raw, &req)
	if err != nil {
		encoded, _ := json.Marshal(&MessageResponse{
			Error: &models.Error{
				Code:    "json",
				Message: "decoding error",
			},
		})
		br.Broadcast(encoded)
	}

	switch req.Type {
	case "command":
		extra, _ := json.Marshal(&CommandExtra{
			UID: req.UID,
			RID: req.RID,
		})
		result := processCommand(&req, extra, cp)
		if result != nil {
			r, _ := json.Marshal(result)
			br.Broadcast(r)
			addMessage(mdb, req.RID, result.Message, "bot")
		}
	case "message":
		res := &MessageResponse{
			Type:    "message",
			RID:     req.RID,
			UID:     req.UID,
			Message: req.Message,
		}
		r, _ := json.Marshal(res)
		br.Broadcast(r)
		addMessage(mdb, req.RID, req.Message, req.UID)
	}
}

func processCommand(req *MessageRequest, extra []byte, cp cmds.CommandProcesor) *MessageResponse {
	if cmd, ok := cp.IsCommand(req.Message); ok {
		cp.Process(&cmds.Request{
			Command: cmd,
			Extra:   string(extra),
		})
		return nil
	}
	return &MessageResponse{
		Error: &models.Error{
			Code:    "command",
			Message: "Not found",
		},
	}
}

func addMessage(mdb chatroom.Message, rid, msg, uid string) {
	mdb.Add(context.TODO(), &models.Message{
		ID:      uuid.New().String(),
		RID:     rid,
		UID:     uid,
		Created: time.Now(),
		Message: msg,
	})
}
