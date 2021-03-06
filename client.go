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
	send        chan *response
	numMessages int32
	muted       bool
}

func (c *client) sendError(err string) {
	resp := &response{
		Protocol: "err",
		Data:     err,
	}

	c.send <- resp
}

func (c *client) readPump() {

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessage)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, bytes, err := c.conn.ReadMessage()

		if err != nil {
			break
		}

		req := &request{}

		err = json.Unmarshal(bytes, req)

		if err != nil {
			c.sendError("Invalid json")
			continue
		}

		req.Client = c

		c.hub.message <- req
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case response, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			bytes, err := json.Marshal(response)
			if err != nil {
				c.sendError("Internal server error")
				log.Println(err)
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
