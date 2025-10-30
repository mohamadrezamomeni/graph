package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	contactHandler "github.com/mohamadrezamomeni/graph/delivery/httpserver/controller/contact"
	appError "github.com/mohamadrezamomeni/graph/pkg/error"
	appLogger "github.com/mohamadrezamomeni/graph/pkg/log"
	contactService "github.com/mohamadrezamomeni/graph/service/contact"
	contactValidator "github.com/mohamadrezamomeni/graph/validator/contact"
)

type Handler interface {
	SetRouter(*echo.Group)
}

type Server struct {
	router         *echo.Echo
	httpConfig     *HTTPConfig
	contactHandler Handler
}

func New(
	cfg *HTTPConfig,
	contactSvc *contactService.Contact,
	contactValidator *contactValidator.Validator,
) *Server {
	return &Server{
		router: echo.New(),
		contactHandler: contactHandler.New(
			contactSvc,
			contactValidator,
		),
		httpConfig: cfg,
	}
}

func (s *Server) Serve() {
	scope := "httpserver.Serve"

	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Recover())
	s.router.Logger.SetOutput(appLogger.Writer())
	s.router.HideBanner = true

	api := s.router.Group("/api/v1")

	s.contactHandler.SetRouter(api)

	address := fmt.Sprintf(":%s", s.httpConfig.Port)
	err := s.router.Start(address)
	if err != nil && err != http.ErrServerClosed {
		panic(
			appError.
				Wrap(err).
				Scope(scope).
				DeactiveWrite().
				Errorf("error to start httpserver"),
		)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
