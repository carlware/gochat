package cmds

type stock struct{}

func (s *stock) Execute(arg string) (string, error) {
	return "", nil
}

func (s *stock) Prepare(req *Request) string {
	return req.Command.Argument
}

func (s *stock) Type() string {
	return "queue"
}
