package cmds

import (
	"encoding/json"
	"strings"

	"github.com/carlware/gochat/chatroom/models"
	"github.com/carlware/gochat/common/mq"
)

var commands CommandHash

type listener struct {
	commandSender   mq.Sender
	commandReceiver mq.Listener
	send            chan []byte
	results         chan *Response
}

func init() {
	commands = CommandHash{}
}

func NewCommandProcessor(sender mq.Sender, receiver mq.Listener) *listener {
	return &listener{
		send:            make(chan []byte),
		results:         make(chan *Response),
		commandSender:   sender,
		commandReceiver: receiver,
	}
}

func (p *listener) Process(req *Request) {
	if cmd, ok := commands[req.Command.Name]; ok {
		switch cmd.Type() {
		case "queue":
			prep := cmd.Prepare(req)
			raw, _ := json.Marshal(&MQRequest{
				Command: prep,
				Extra:   req.Extra,
			})
			p.commandSender.Send(raw)
		case "fast":
			result := cmd.Execute(req.Command.Name)
			p.results <- &Response{
				Result: result,
			}
		}
	}
}

func (p *listener) Results() chan *Response {
	return p.results
}

func (p *listener) IsCommand(raw string) (*Command, bool) {
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

func (p *listener) Run() {

	// Listen for MQ broker results
	go func() {
		consumer, _ := p.commandReceiver.Listen()
		for msg := range consumer {
			mRes := &MQResponse{}
			res := &Response{}
			err := json.Unmarshal(msg, mRes)
			if err != nil {
				res.Error = &models.Error{
					Code:    "json",
					Message: "cannot decode mq result",
				}
			} else {
				res.Result = mRes.Result
				res.Extra = mRes.Extra
				res.Error = mRes.Error
			}
			p.results <- res
		}
	}()
}
