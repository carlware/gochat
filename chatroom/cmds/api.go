package cmds

type Command struct {
	Name     string `json:"name"`
	Argument string `json:"argument"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Request struct {
	Extra   string   `json:"extra"`
	Command *Command `json:"command"`
}

type Response struct {
	Response string `json:"response"`
	Extra    string `json:"extra"`
	Error    Error  `json:"error"`
}

type Executor interface {
	Execute(arg string) *Response
	Prepare(*Request) *Request
	Type() string
}

type CommandHash map[string]Executor
