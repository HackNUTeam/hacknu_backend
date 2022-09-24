package handler

import (
	"hacknu/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	services       *service.Service
	upgrader       *websocket.Upgrader
	dispatcherChan chan []byte
}

func NewHandler(services *service.Service, upgrader websocket.Upgrader, dispChan chan []byte) *Handler {
	return &Handler{services: services, upgrader: &upgrader, dispatcherChan: dispChan}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/", h.serveHome)
	router.GET("/user", h.HandleUser)
	router.GET("/dispatcher", h.HandleDispatcher)
	router.GET("/get-history", h.GetHistory)
	return router
}
