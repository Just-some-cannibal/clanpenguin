package main

import (
	"encoding/json"
	"log"
	"time"
)

const (
	reallocnum  = 100
	maxMessages = 5
)

//Hub contains all necessary info for the server
type Hub struct {
	clients    map[*client]bool
	message    chan *request
	register   chan *client
	unregister chan *client
	messages   []message
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		message:    make(chan *request),
		messages:   make([]message, 0, reallocnum),
	}
}

func (h *Hub) reallocMessages(size int) {
	temp := make([]message, 0, size)
	copy(h.messages, temp)
	h.messages = temp
}

func (h *Hub) pushMessage(m message) {
	if cap(h.messages) < len(h.messages)+1 {
		h.reallocMessages(cap(h.messages) + reallocnum)
	}
	h.messages = append(h.messages, m)
}

func (h *Hub) broadcast(req *request) {
	c := req.Client
	if c.numMessages == maxMessages {
		c.sendError("spam")

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

	if c.numMessages == 0 {
		go func() {
			timer := time.NewTimer(5 * time.Second)
			<-timer.C
			c.numMessages = 0
		}()
	}

	c.numMessages++

	var m = message{}

	err := json.Unmarshal(req.Data, &m)

	if err != nil {
		req.Client.sendError("internal")
		return
	}

	if len(m.Text) > 100 || len(m.User) > 20 {
		c.sendError("maxlength")
		return
	}

	h.pushMessage(m)

	resp := &response{
		Data:     m,
		Protocol: "broadcast",
	}

	for client := range h.clients {
		select {
		case client.send <- resp:
		default:
			delete(h.clients, client)
			close(client.send)
		}
	}
}

func (h *Hub) get(r *request) {
	var messages []message

	if len(h.messages) > 100 {
		messages = h.messages[len(h.messages)-100:]
	} else {
		messages = h.messages[:]
	}

	resp := &response{
		Data:     messages,
		Protocol: "get",
	}

	r.Client.send <- resp
}

//Run is a goroutine that has a handler for all of its channels
func (h *Hub) Run() {
	log.Println("Running hub")
	for {
		select {
		case client := <-h.register:
			log.Println("Registering ", client.conn.RemoteAddr().String())
			h.clients[client] = true
		case client := <-h.unregister:
			log.Println("Unregistering ", client.conn.RemoteAddr().String())
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
