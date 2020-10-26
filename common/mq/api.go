package mq

type Sender interface {
	Send([]byte) error
}

type Listener interface {
	Listen() ([]byte, error)
}
