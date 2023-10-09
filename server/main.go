package main

import (
	"log"
	"net/http"
)

func main() {
	setupApi()

	log.Fatal(http.ListenAndServe(port, nil))
}

func setupApi() {
	manager := NewManager()

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/static", http.FileServer(http.Dir("static/static")))
	http.HandleFunc("/ws", manager.serveWebsocket)
}
