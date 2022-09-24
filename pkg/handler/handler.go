package handler

import (
	"hacknu/model"
	"hacknu/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	services *service.Service
	hub      *model.Hub
	upgrader *websocket.Upgrader
}

func NewHandler(services *service.Service, hub *model.Hub, upgrader websocket.Upgrader) *Handler {
	return &Handler{services: services, hub: hub, upgrader: &upgrader}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/", h.serveHome)
	router.GET("/ws", h.ServeWs)
	//router.POST("/sendLocation", h.sendLocation)
	return router
}
