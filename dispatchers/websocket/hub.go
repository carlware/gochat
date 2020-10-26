package websocket

import (
	"fmt"

	log "github.com/inconshreveable/log15"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Output messages to external interfaces (dependency injection)
	send chan []byte
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		send:       make(chan []byte),
	}
}

func (h *Hub) Run() {
	log.Info("web socket start to running")
	for {
		select {
		case client := <-h.register:
			fmt.Println("new client registered")
			h.clients[client] = true
		case client := <-h.unregister:
			fmt.Println("new client unregistered")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			fmt.Println("message to client", message)
			h.send <- message
			// stockbot.SendMessage(message)
			// for client := range h.clients {
			// 	select {
			// 	case client.send <- message:
			// 	default:
			// 		close(client.send)
			// 		delete(h.clients, client)
			// 	}
			// }
		}
	}
}

func (h *Hub) Receive() (chan []byte, error) {
	return h.send, nil
}

func (h *Hub) Broadcast(msg []byte) {
	for client := range h.clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
