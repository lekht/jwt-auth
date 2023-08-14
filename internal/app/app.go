package app

import (
	"context"
	"fmt"
	"jwt-auth/config"
	"jwt-auth/internal/auth"
	"jwt-auth/internal/controllers"
	"jwt-auth/internal/repository"
	"jwt-auth/package/httpserver"
	"jwt-auth/package/postgres"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	db, err := postgres.New(context.Background(), &cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepo(db)
	if err != nil {
		log.Fatal("NewStorage:", err)
	}

	authentificator := auth.New(cfg.Secret, repo)

	// HTTP Server
	handler := gin.New()

	controllers.NewRouter(handler, authentificator)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.Server.Port))

	// Ожидает interrupt сигнал
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
