package cases

import (
	"encoding/json"
	"fmt"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/cmds"
	"github.com/carlware/gochat/chatroom/models"
)

type mServer struct {
	br chatroom.BroadcastReceiver
	cp cmds.CommandProcesor
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

func NewMessageListener(br chatroom.BroadcastReceiver, cp cmds.CommandProcesor) *mServer {
	return &mServer{
		br: br,
		cp: cp,
	}
}

func (s *mServer) Listen() {
	// listen
	fmt.Println("listen incoming messages started")
	messages, _ := s.br.Receive()

	// Listen incoming messages
	go func() {
		for msg := range messages {
			fmt.Println("incoming message", string(msg))
			messageProccesor(msg, s.br, s.cp)
		}
	}()

	// Listen incoming responses from cmds manager
	go func() {
		results := s.cp.Results()
		for result := range results {
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
		}
	}()

	fmt.Println("incomming message listener close")
}

func messageProccesor(raw []byte, br chatroom.BroadcastReceiver, cp cmds.CommandProcesor) {
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
		}
	case "message":
		res := &MessageResponse{
			Type:    "message",
			RID:     req.RID,
			UID:     req.UID,
			Message: string(raw),
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
