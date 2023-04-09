package mq

type Sender interface {
	Send([]byte, string) error
}

type Listener interface {
	Listen(string) (chan []byte, error)
}

type ListenSender interface {
	Sender
	Listener
}
