package stock

import "github.com/carlware/gochat/chatroom/cmds"

type stock struct{}

func init() {
	cmds.Add("stock", &stock{})
}

func (s *stock) Execute(arg string) *cmds.Response {
	return nil
}

func (s *stock) Prepare(req *cmds.Request) *cmds.Request {
	return req
}

func (s *stock) Type() string {
	return "queue"
}
