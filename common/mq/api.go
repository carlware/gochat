package mq

type Sender interface {
	Send([]byte) error
}

type Listener interface {
	Listen() (chan []byte, error)
}

type ListenSender interface {
	Sender
	Listener
}
