package config

import (
	"github.com/mohamadrezamomeni/graph/delivery/httpserver"
	"github.com/mohamadrezamomeni/graph/repository/sqlite"
)

func Load() *Config {
	return &Config{
		Application: Application{
			GracefulShutdownTimeout: 20,
		},
		HTTP: httpserver.HTTPConfig{
			Port: "1234",
		},
		DB: sqlite.DBConfig{
			Path: "graph.db",
		},
	}
}
