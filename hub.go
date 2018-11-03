package main

import (
	"encoding/json"
	"log"
	"time"
)

const (
	reallocnum = 100
)

//Hub contains all necessary info for the server
type Hub struct {
	clients    map[*client]bool
	message    chan *Request
	register   chan *client
	unregister chan *client
	messages   []Message
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		message:    make(chan *Request),
		messages:   make([]Message, 0, reallocnum),
	}
}

func (h *Hub) reallocMessages(size int) {
	temp := make([]Message, 0, size)
	copy(h.messages, temp)
	h.messages = temp
}

func (h *Hub) pushMessage(message Message) {
	if cap(h.messages) < len(h.messages)+1 {
		h.reallocMessages(cap(h.messages) + reallocnum)
	}
	h.messages = append(h.messages, message)
}

func (h *Hub) broadcast(request *Request) {
	c := request.Client
	if c.numMessages == 5 {
		c.sendError("Too many messages")

		if !c.muted {
			c.muted = true

			go func() {
				timer := time.NewTimer(5 * time.Second)
				<-timer.C
				c.muted = false
				c.numMessages = 0
			}()
		}
		return
	}

	c.numMessages++

	var message = Message{}

	err := json.Unmarshal(request.Data, &message)

	if err != nil {
		request.Client.sendError("Internal server error")
		return
	}

	if len(message.Text) > 100 || len(message.User) > 20 {
		c.sendError("Exceeded max message length")
		return
	}

	h.pushMessage(message)

	response := &Response{
		Data:     message,
		Protocol: "broadcast",
	}

	for client := range h.clients {
		select {
		case client.send <- response:
		default:
			delete(h.clients, client)
			close(client.send)
		}
	}
}

func (h *Hub) get(r *Request) {
	var messages []Message

	if len(h.messages) > 100 {
		messages = h.messages[len(h.messages)-100:]
	} else {
		messages = h.messages[:]
	}

	response := &Response{
		Data:     messages,
		Protocol: "get",
	}

	r.Client.send <- response
}

//Run is a goroutine that has a handler for all of its channels
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

		case m := <-h.message:
			switch m.Protocol {
			case "get":
				h.get(m)
			case "broadcast":
				h.broadcast(m)
			default:
				log.Println("Invalid request ", m.Protocol)
			}
		}
	}
}
