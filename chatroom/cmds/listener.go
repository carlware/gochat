package cmds

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/carlware/gochat/common/mq"
)

type listener struct {
	commands CommandHash
}

var commands CommandHash
var mQ *mq.ListenSender
var send chan []byte
var results chan *Response

func init() {
	send = make(chan []byte)
	results = make(chan []byte)
	commands = CommandHash{}
}

func Process(req *Request) error {
	if cmd, ok := commands[req.Command.Name]; ok {
		switch cmd.Type() {
		case "queue":
			req := cmd.Prepare(req)
			if mQ != nil {
				raw, err := json.Marshal(req)
				if err != nil {
					return err
				}
				mQ.Send(raw)
			} else {
				return errors.New("MQ broker is not set")
			}
		case "fast":
			result := cmd.Execute(req.Command.Name)
			results <- &result
		}
	}
}

func IsCommand(raw string) (*Command, bool) {
	token := strings.Split(raw, "=")
	if len(token) != 2 {
		return nil, false
	}
	return &Command{
		Name:     token[0],
		Argument: token[1],
	}, true
}

func Add(name string, command Executor) {
	commands[name] = command
}

func SetMQ(mq *mq.ListenSender) {
	if mQ != nil {
		mQ = mq
	}
}

func Results() chan *Response {
	return results
}

func Listen() {
	go func() {
		consumer := mQ.Listen()
		for msg := range consumer {
			res := &Response{}
			err := json.Unmarshal(msg, res)
			if err != nil {
				res.Error = Error{
					Code:    "json",
					Message: "Cannot decode",
				}
			}
			results <- Response
		}
	}()
}
