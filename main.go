package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	server *Server
)

func main() {
	initializeConfig()

	port := viper.GetString("port")

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

func initializeConfig() {
	viper.SetDefault("port", 8081)

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No configuration file found, using defaults.")
	}
}
