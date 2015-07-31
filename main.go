package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	server *Server
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	server = NewServer()

	http.Handle("/", http.FileServer(http.Dir("./client/")))
	http.HandleFunc("/chat", chatHandler)

	log.Printf("Starting server on port '%s'", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func readLoop(conn *websocket.Conn) {
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request for this endpoint", 405)
	}

	ws, err := upgrader.Upgrade(w, r, http.Header{
		"Sec-websocket-Protocol": websocket.Subprotocols(r),
	})
	if err != nil {
		log.Println("Error Upgrading: ", err.Error())
		return
	}

	cli := NewClient(server, ws)

	cli.Handle()
}
