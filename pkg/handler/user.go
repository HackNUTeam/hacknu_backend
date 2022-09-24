package handler

import (
	"hacknu/model"
	"log"
	"net/http"
	"os"

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

func (h *Handler) ServeWs(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	log.Print(conn)
	if err != nil {
		log.Println(err)
		return
	}
	client := &model.Client{Hub: h.hub, Conn: conn, Send: make(chan *model.LocationData, 256)}
	log.Print(client)
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

/*
func (h *Handler) sendLocation(c *gin.Context) {
	var location model.LocationData
	if err := c.BindJSON(&location); err != nil {
		defaultErrorHandler(c, errors.New("bad request | "+err.Error()))
		return
	}

	go client.WritePump()
}
*/
