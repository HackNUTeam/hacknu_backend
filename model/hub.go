package model

import (
	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}
