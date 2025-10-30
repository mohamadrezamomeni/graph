package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mohamadrezamomeni/graph/delivery/httpserver"
	"github.com/mohamadrezamomeni/graph/pkg/config"
	appLogger "github.com/mohamadrezamomeni/graph/pkg/log"
	"github.com/mohamadrezamomeni/graph/repository/migrate"
)

func main() {
	cfg := config.Load()

	migration := migrate.New(&cfg.DB)

	migration.UP()

	server := httpserver.New(
		&cfg.HTTP,
	)

	go func() {
		server.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxWithTimeout); err != nil {
		appLogger.Info(err.Error())
	}
}
