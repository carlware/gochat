package cases

import (
	"context"
	"encoding/json"
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
	Created time.Time     `json:"created"`
	ID      string        `json:"id"`
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
			extra := CommandExtra{}
			_ = json.Unmarshal([]byte(result.Extra), &extra)
			message := result.Result
			if result.Error != nil {
				message = result.Error.Message
			}
			m := addMessage(s.messagedb, extra.RID, message, "bot")
			response := &MessageResponse{
				Message: result.Result,
				Type:    "command",
				RID:     extra.RID,
				UID:     m.UID,
				Error:   result.Error,
				Created: m.Created,
				ID:      m.ID,
			}
			encoded, _ := json.Marshal(response)
			s.br.Broadcast(encoded)
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
			m := addMessage(mdb, req.RID, result.Message, "bot")
			result.Created = m.Created
			result.ID = m.ID
			result.UID = m.UID
			r, _ := json.Marshal(result)
			br.Broadcast(r)
		}
	case "message":
		m := addMessage(mdb, req.RID, req.Message, req.UID)
		res := &MessageResponse{
			Type:    "message",
			RID:     req.RID,
			UID:     req.UID,
			Message: req.Message,
			Created: m.Created,
			ID:      m.ID,
		}
		r, _ := json.Marshal(res)
		br.Broadcast(r)
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

func addMessage(mdb chatroom.Message, rid, msg, uid string) *models.Message {
	m := &models.Message{
		ID:      uuid.New().String(),
		RID:     rid,
		UID:     uid,
		Created: time.Now(),
		Message: msg,
	}
	mdb.Add(context.TODO(), m)
	return m
}
