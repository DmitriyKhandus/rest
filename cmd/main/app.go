package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/DmitriyKhandus/rest-api/internal/config"
	"github.com/DmitriyKhandus/rest-api/internal/user"
	"github.com/DmitriyKhandus/rest-api/internal/user/db"
	"github.com/DmitriyKhandus/rest-api/pkg/client/mongodb"
	"github.com/DmitriyKhandus/rest-api/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()
	cfg := config.GetConfig()
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfg.MongoDB.Host,
		cfg.MongoDB.Port, cfg.MongoDB.Username, cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB)
	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	user1 := user.User{
		Id:           "",
		Email:        "myemail@mail.ru",
		Username:     "imagination",
		PasswordHash: "123",
	}

	user1Id, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(user1Id)

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("create socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("Server is listening on socket %s", socketPath)

	} else {
		logger.Info("create tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("Server is listening on port %s", cfg.Listen.Port)
	}
	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	server.Serve(listener)
}
