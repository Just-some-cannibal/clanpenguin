package main

import (
	"log"
)

type Hub struct {
	clients    map[*client]bool
	get        chan *Message
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		get:        make(chan *Message),
	}
}

func (h *Hub) Run() {
	log.Println("Running hub")
	for {
		select {
		case client := <-h.register:
			log.Println("Registering socket")
			h.clients[client] = true
		case client := <-h.unregister:
			log.Println("Unregistering socket")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println("Sending all clients message")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		case m := <-h.get:
			switch m.Protocol {
			case "get":
				m.Client.send <- []byte(m.Data)
			case "send":
			}
		}
	}
}
