package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/carlware/gochat/dispatchers/websocket"
	"github.com/rs/cors"
)

var addr = flag.String("addr", ":8080", "http server address")
var ctx = context.Background()

func GetCorsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"},
	})
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	hub := websocket.NewHub()
	go hub.Run()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	fs := http.FileServer(http.Dir("./web/build"))
	mux.Handle("/", fs)

	handler := GetCorsConfig().Handler(mux)
	log.Fatal(http.ListenAndServe(*addr, handler))
	fmt.Println("starting")
}
