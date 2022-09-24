package handler

import (
	"hacknu/model"
	"hacknu/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	services   *service.Service
	hub        *model.Hub
	upgrader   *websocket.Upgrader
	dispatcher *model.Client
}

func NewHandler(services *service.Service, hub *model.Hub, upgrader websocket.Upgrader) *Handler {
	return &Handler{services: services, hub: hub, upgrader: &upgrader}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/", h.serveHome)
	router.GET("/ws", h.ServeWs)
	router.GET("/ws-disp", h.SendLocation)
	//router.POST("/sendLocation", h.sendLocation)
	return router
}
