package handler

import (
	"bytes"
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
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
	//go client.WritePump()
	go h.ReadPump(ctx, client)
}

func (h *Handler) SendLocation(c *gin.Context) {
	pong := &model.PongStruct{}
	h.pong <- pong
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	var messageByte []byte
	select {
	case <-ctx.Done():
		log.Println("pong didn't received, no clinet")
		return
	case messageByte = <-h.ping:
		log.Println("pong received: starting to proccess messages")
	}
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	log.Print(conn)
	if err != nil {
		log.Println(err)
		return
	}
	client := &model.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}

	go h.WritePump(client, messageByte)
	close(h.ping)
}

func (h *Handler) ReadPump(ctx context.Context, c *model.Client) {
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
		select {
		case <-ctx.Done():
			log.Println("ebu ernara")
			return
		case <-h.pong:
			log.Println("pong received: dispetcher alive")
		}
		h.ping <- message
		var location model.LocationData
		_ = json.Unmarshal(message, &location)
		log.Print(location)
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Print(message)
		c.Hub.Broadcast <- &location
	}
}

func (h *Handler) WritePump(c *model.Client, msg []byte) {
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
				c.Conn.WriteMessage(websocket.CloseMessage, msg)
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
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("BRUH %v", err)
				return
			}
		}
	}
}
