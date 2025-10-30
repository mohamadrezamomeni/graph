package config

import (
	"time"

	"github.com/mohamadrezamomeni/graph/delivery/httpserver"
	"github.com/mohamadrezamomeni/graph/repository/sqlite"
)

type Application struct {
	GracefulShutdownTimeout time.Duration
}

type Config struct {
	Application Application
	HTTP        httpserver.HTTPConfig
	DB          sqlite.DBConfig
}
