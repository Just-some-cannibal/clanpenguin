package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessage = 512
)

type client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan *Response
	numMessages int32
	muted       bool
}

func (c *client) sendError(err string) {
	response := &Response{
		Protocol: "err",
		Data:     err,
	}

	c.send <- response
}

func (c *client) readPump() {
	log.Println("Registering Read Pump")

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessage)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, bytes, err := c.conn.ReadMessage()

		log.Println("Received message")

		if err != nil {
			log.Println("Error receiving message\n", err)
			break
		}

		message := &Request{}

		err = json.Unmarshal(bytes, message)

		if err != nil {
			log.Println("Invalid format")
			continue
		}

		message.Client = c

		c.hub.message <- message
	}
}

func (c *client) writePump() {
	log.Println("Initializing write pump")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case response, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			log.Println("Sending data")
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("Error getting writer\n", err)
				return
			}

			bytes, err := json.Marshal(response)
			if err != nil {
				log.Println("Could not marshal the data")
				continue
			}

			w.Write(bytes)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}
