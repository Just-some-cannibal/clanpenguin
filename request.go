package main

import (
	"encoding/json"
)

// message is a struct that holds the user message
type message struct {
	Text string `json:"text"`
	User string `json:"user"`
}

//request is a struct that has all information used for a request
type request struct {
	Protocol string          `json:"protocol"`
	Data     json.RawMessage `json:"data"`
	Auth     string          `json:"auth"`
	Client   *client         `json:"-"`
}

//response is a struct that gives the status code and associated data
type response struct {
	Protocol string      `json:"protocol"`
	Data     interface{} `json:"data"`
}
