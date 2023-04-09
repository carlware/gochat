package cmds

import (
	"encoding/json"
	"strings"

	"github.com/carlware/gochat/chatroom/models"
	"github.com/carlware/gochat/common/mq"
	log "github.com/inconshreveable/log15"
)

var commands CommandHash

type listener struct {
	mq      mq.ListenSender
	send    chan []byte
	results chan *Response
}

func init() {
	commands = CommandHash{}
	commands["stock"] = &stock{}
}

func NewCommandProcessor(ls mq.ListenSender) *listener {
	return &listener{
		send:    make(chan []byte),
		results: make(chan *Response),
		mq:      ls,
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
			p.mq.Send(raw, "request")
		case "fast":
			result, err := cmd.Execute(req.Command.Name)
			resp := &Response{
				Result: result,
			}
			if err != nil {
				resp.Error = &models.Error{
					Code:    "command",
					Message: err.Error(),
				}
			}
			p.results <- resp
		default:
			log.Info("cannot process command", "type", cmd.Type())
		}
	} else {
		p.results <- &Response{
			Extra: req.Extra,
			Error: &models.Error{
				Code:    "command",
				Message: "command does not exist",
			},
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

func (p *listener) Run() {
	// Listen for MQ broker results
	go func() {
		consumer, _ := p.mq.Listen("results")

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
