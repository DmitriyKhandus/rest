package main

import (
	"net"
	"net/http"
	"time"

	"github.com/DmitriyKhandus/rest-api/internal/user"
	"github.com/DmitriyKhandus/rest-api/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()
	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	logger.Info("Server is listening on port: 3000")
	server.Serve(listener)
}
