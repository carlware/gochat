package cases

import (
	"fmt"

	"github.com/carlware/gochat/chatroom"
	"github.com/carlware/gochat/chatroom/cmds"
)

type Options struct {
	BroadcastReceiver chatroom.BroadcastReceiver
}

func ListenMessages(opts *Options) {
	// listen
	fmt.Println("listen incoming messages started")
	messages, _ := opts.BroadcastReceiver.Receive()

	// Listen incoming messages
	go func() {
		for msg := range messages {
			fmt.Println("incoming message", string(msg))
			raw := string(msg)
			if cmd, ok := cmds.IsCommand(); ok {

			}

			opts.BroadcastReceiver.Broadcast(msg)
		}
	}()

	fmt.Println("incomming message listener close")
}
