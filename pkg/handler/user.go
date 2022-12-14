package handler

import (
	"encoding/json"
	"errors"
	"hacknu/model"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
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
	h.dispatcherChan = make(chan []byte, 1000)

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
			var location model.LocationData
			err := json.Unmarshal(msg, &location)
			if err != nil {
				log.Print(err)
			}
			h.dispatcherChan <- msg
			err = h.services.User.CreateReading(&location)
			if err != nil {
				log.Print(err)
			}
		} else {
			log.Printf("Dispatcher channel is nil")
		}
	}
}

func (h *Handler) GetHistory(c *gin.Context) {
	req := &model.GetLocationRequest{
		Timestamp: -1,
	}
	req.Name = c.Query("name")
	var err error
	log.Print(req.Name)
	keyTime := c.Query("timestamp")
	if keyTime != "" {
		req.Timestamp, err = strconv.ParseInt(keyTime, 10, 64)
	}

	log.Print(req.Timestamp)
	if err != nil {
		panic(err)
	}
	res, err := h.services.User.GetHistory(req)
	if err != nil {
		log.Printf("get locations error: %v", err)
		if errors.Is(err, model.ErrNoDataForSuchUser) {
			c.AbortWithStatusJSON(500, createResponse(nil, model.ErrNoDataForSuchUser.Error()))
			return
		}
		c.AbortWithStatusJSON(500, createResponse(nil, "INTERNAL_SERVER_ERROR"))
		return
	}
	c.JSON(200, createResponse(res, ""))
}
func createResponse(data interface{}, err string) gin.H {
	return gin.H{
		"data":  data,
		"error": err,
	}
}
