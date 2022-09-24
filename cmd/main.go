package main

import (
	"errors"
	maps "hacknu"
	"net/http"
	"os"

	"hacknu/pkg/handler"
	"hacknu/pkg/repository"
	"hacknu/pkg/service"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Print("No .env file found, please set")
	}
}

func main() {
	logrus.Print("Startup server")
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initEnv(); err != nil {
		logrus.Fatalf("error initializing env: %s", err.Error())
	}

	db, err := repository.NewPostgreDB(
		os.Getenv("DSN"),
	)

	if err != nil {
		logrus.Fatalf(err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	handlers := handler.NewHandler(service, upgrader, nil)
	staticHandler := handler.NewStaticHandler(service)

	srv := new(maps.Server)
	staticSrv := new(maps.Server)
	go func() {
		if err := srv.Run(os.Getenv("APIPortHTTP"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	if err := staticSrv.Run(os.Getenv("StaticPortHTTP"), staticHandler.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initEnv() error {

	reqs := []string{
		"StaticPortHTTP",
		"APIPortHTTP",
		"Data",
	}

	for i := 0; i < len(reqs); i++ {
		_, exists := os.LookupEnv(reqs[i])

		if !exists {
			return errors.New(".env variables not set")
		}
	}

	return nil
}
