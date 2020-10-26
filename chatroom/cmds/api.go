package cmds

import "github.com/carlware/gochat/chatroom/models"

type Command struct {
	Name     string `json:"name"`
	Argument string `json:"argument"`
}

type Request struct {
	Extra   string   `json:"extra"`
	Command *Command `json:"command"`
}

type Response struct {
	Result string        `json:"response"`
	Extra  string        `json:"extra"`
	Error  *models.Error `json:"error"`
}

type MQRequest struct {
	Command string `json:"command"`
	Extra   string `json:"extra"`
}

type MQResponse struct {
	Result string        `json:"result"`
	Extra  string        `json:"extra"`
	Error  *models.Error `json:"error"`
}

type Executor interface {
	Execute(arg string) string
	Prepare(*Request) string
	Type() string
}

type CommandHash map[string]Executor

type CommandProcesor interface {
	IsCommand(raw string) (*Command, bool)
	Results() chan *Response
	Process(req *Request)
}
