package handler

import (
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

func (h *Handler) HandleUser(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	log.Print(conn)
	if err != nil {
		log.Println(err)
		return
	}
	client := &model.Client{Conn: conn, Send: make(chan []byte, 256)}

	h.listenUser(client)
	log.Printf("Finished listening user %v", client)
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

func readDispChan(errChan chan error, client *model.Client) {
	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			errChan <- err
			log.Println("Error during message reading:", err)
			break
		}
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
		} else {
			log.Printf("Dispatcher channel is nil")
		}
	}
}
