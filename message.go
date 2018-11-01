package main

import (
	_ "encoding/json"
)

type Message struct {
	Protocol string  `json:"message"`
	Data     string  `json:"data"`
	Auth     string  `json:"auth"`
	Client   *client `json:"-"`
}
