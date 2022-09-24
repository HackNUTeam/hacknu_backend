package handler

import (
	"encoding/json"
	"errors"
	"hacknu/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

<<<<<<< HEAD
func (h *Handler) ServeWs(c *gin.Context) {
	//h.ping = make(chan []byte, 256)
=======
func (h *Handler) HandleUser(c *gin.Context) {
>>>>>>> origin/new-socket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	log.Print(conn)
	if err != nil {
		log.Println(err)
		return
	}
<<<<<<< HEAD
	client := &model.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}
	log.Print(client)
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go h.WritePump(client)
	go h.ReadPump(client)
=======
	client := &model.Client{Conn: conn, Send: make(chan []byte, 256)}

	h.listenUser(client)
	log.Printf("Finished listening user %v", client)
>>>>>>> origin/new-socket
}

func (h *Handler) HandleDispatcher(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	log.Print(conn)
	if err != nil {
		log.Println(err)
		return
	}
	client := &model.Client{Conn: conn, Send: make(chan []byte, 256)}
	h.dispatcherChan = make(chan []byte, 100)

	h.listenDispatcherChan(client)
}

func (h *Handler) listenDispatcherChan(client *model.Client) {
	errChan := make(chan error, 10)
	go readDispChan(errChan, client)
	for {

		select {
		case msg := <-h.dispatcherChan:
			log.Printf("Received message for dispatcher %v", msg)
			client.Conn.WriteMessage(1, msg)
		case err := <-errChan:
			log.Printf("Lost connection to dispatcher: %v", err)
			h.dispatcherChan = nil
			return
		}
	}
}

<<<<<<< HEAD
func (h *Handler) ReadPump(c *model.Client) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
=======
func readDispChan(errChan chan error, client *model.Client) {
>>>>>>> origin/new-socket
	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			errChan <- err
			log.Println("Error during message reading:", err)
			break
		}
<<<<<<< HEAD
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
=======
	}
}

func (h *Handler) listenUser(client *model.Client) {
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client %v", err)
			return
		}
		log.Printf("Received message from client %v", msg)
		if h.dispatcherChan != nil {
			h.dispatcherChan <- msg
			var location model.LocationData
			err := json.Unmarshal(msg, &location)
>>>>>>> origin/new-socket
			if err != nil {
				log.Print(errors.New("Could not unmarshall"))
				return
			}
			err = h.services.User.CreateReading(&location)
			if err != nil {
				log.Print(err)
				return
			}
		} else {
			log.Printf("Dispatcher channel is nil")
		}
	}
}
