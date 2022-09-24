package handler

import (
	"encoding/json"
	"hacknu/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (h *Handler) serveHome(c *gin.Context) {
	log.Println(c.Request.URL.Path)

	if c.Request.URL.Path != "/" {
		http.Error(c.Writer, "Not found", http.StatusNotFound)
		return
	}
	if c.Request.Method != http.MethodGet {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println(c.Request.URL.Path == "/")

	http.ServeFile(c.Writer, c.Request, os.Getenv("Data")+"home.html")
}

func (h *Handler) ServeWs(c *gin.Context) {
	//h.ping = make(chan []byte, 256)
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	log.Print(conn)
	if err != nil {
		log.Println(err)
		return
	}
	client := &model.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}
	log.Print(client)
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go h.WritePump(client)
	go h.ReadPump(client)
}

func (h *Handler) SendLocation(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	log.Print(conn)
	if err != nil {
		log.Println(err)
		return
	}
	client := &model.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}
	h.dispatcher = client
	//h.pong = make(chan *model.PongStruct, 10)
}

func (h *Handler) ReadPump(c *model.Client) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		log.Print(message, err)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		h.WritePump(h.dispatcher)
	}
}

func (h *Handler) WritePump(c *model.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Print("Hub closed by server")
				c.Conn.WriteMessage(websocket.CloseMessage, message)
				return
			}
			log.Print(message)
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print(err)
				return
			}
			send, err := json.Marshal(message)
			if err != nil {
				log.Print(err)
				return
			}
			w.Write(send)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				read := <-c.Send
				newSend, err := json.Marshal(read)
				if err != nil {
					log.Print(err)
					return
				}
				w.Write(newSend)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
