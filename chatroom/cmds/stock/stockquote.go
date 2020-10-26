package stock

import "github.com/carlware/gochat/chatroom/cmds"

type stock struct{}

func init() {
	cmds.Add("stock", &stock{})
}

func (s *stock) Execute(arg string) string {
	return ""
}

func (s *stock) Prepare(req *cmds.Request) string {
	return req.Command.Name + "=" + req.Command.Argument
}

func (s *stock) Type() string {
	return "queue"
}
