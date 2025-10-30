package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/mohamadrezamomeni/graph/delivery/httpserver"
	"github.com/mohamadrezamomeni/graph/pkg/config"
	appLogger "github.com/mohamadrezamomeni/graph/pkg/log"
	"github.com/mohamadrezamomeni/graph/repository/migrate"
	"github.com/mohamadrezamomeni/graph/repository/sqlite"
	contactSqlite "github.com/mohamadrezamomeni/graph/repository/sqlite/contact"
	contactService "github.com/mohamadrezamomeni/graph/service/contact"
	contactValidator "github.com/mohamadrezamomeni/graph/validator/contact"
)

func main() {
	cfg := config.Load()

	migration := migrate.New(&cfg.DB)

	migration.UP()

	contactSvc, contactValidator := setupServices(cfg)

	server := httpserver.New(
		&cfg.HTTP,
		contactSvc,
		contactValidator,
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

func setupServices(cfg *config.Config) (
	*contactService.Contact,
	*contactValidator.Validator,
) {
	db := sqlite.New(&cfg.DB)

	contactRepo := contactSqlite.New(db)
	contactSvc := contactService.New(contactRepo)

	contactValidator := contactValidator.New()
	return contactSvc, contactValidator
}
