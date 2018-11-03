package main

import (
	"encoding/json"
)

// Message is a struct that holds the user message
type Message struct {
	Text string `json:"text"`
	User string `json:"user"`
}

//Request is a struct that has all information used for a request
type Request struct {
	Protocol string          `json:"protocol"`
	Data     json.RawMessage `json:"data"`
	Auth     string          `json:"auth"`
	Client   *client         `json:"-"`
}

//Response is a struct that gives the status code and associated data
type Response struct {
	Protocol string      `json:"protocol"`
	Data     interface{} `json:"data"`
}
