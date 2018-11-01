package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = newHub()

func serveGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving game")
	http.ServeFile(w, r, "game.html")
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving home")
	http.ServeFile(w, r, "home.html")
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating websocket")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("An error has occured registering the socket: ", err)
		return
	}

	client := &client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}

func main() {
	log.Println("Starting server")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	go hub.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/game", serveGame)
	http.HandleFunc("/ws", handleWS)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe("0.0.0.0:8080", nil)
}
