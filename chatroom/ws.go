package chatroom

type Receiver interface {
	Receive() (chan []byte, error)
}

type Broadcaster interface {
	Broadcast([]byte)
}

type BroadcastReceiver interface {
	Receiver
	Broadcaster
}
