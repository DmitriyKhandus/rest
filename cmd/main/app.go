package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/DmitriyKhandus/rest-api/internal/user"
	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Println("create router")
	router := httprouter.New()
	fmt.Println("register user handler")
	handler := user.NewHandler()
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
	fmt.Println("start application")
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	fmt.Println("Server is listening on port: 3000")
	server.Serve(listener)
}
