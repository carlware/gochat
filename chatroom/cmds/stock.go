package cmds

type stock struct{}

func (s *stock) Execute(arg string) (string, error) {
	return "test", nil
}

func (s *stock) Prepare(req *Request) string {
	return req.Command.Name + "=" + req.Command.Argument
}

func (s *stock) Type() string {
	return "queue"
}
