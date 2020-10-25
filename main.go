package main

import (
	"context"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http server address")
var ctx = context.Background()

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir("./web/build"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
