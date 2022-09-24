package handler

import (
	"hacknu/model"
	"hacknu/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	services       *service.Service
	upgrader       *websocket.Upgrader
	clients        map[*model.Client]model.LocationData
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
