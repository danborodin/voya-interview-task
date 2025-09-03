package server

import (
	"context"
	"errors"
	"interview-go/config"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	backendbeer "interview-go/backend/client"
	beerapi "interview-go/internal/beer"
)

type Server struct {
	cfg  *config.Configuration
	Echo *echo.Echo
}

func NewServer(cfg *config.Configuration) (*Server, error) {
	s := &Server{
		cfg: cfg,
	}
	s.newEchoServer()
	return s, nil
}

func (s *Server) StartServer() error {
	server := &http.Server{Addr: ":" + s.cfg.Server.Port}
	server.SetKeepAlivesEnabled(true)

	log.Infof("Starting %s on port %s", s.cfg.App.Name, s.cfg.Server.Port)

	return s.Echo.StartServer(server)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}

func (s *Server) newEchoServer() {
	e := echo.New()
	e.HideBanner = true

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var he *echo.HTTPError
		if errors.As(err, &he) && he.Code > 0 {
			code = he.Code
		}
		resp := map[string]any{
			"error": http.StatusText(code),
			"msg":   err.Error(),
		}
		if !c.Response().Committed {
			_ = c.JSON(code, resp)
		}
	}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			return next(c)
		}
	})

	s.Echo = e

	s.Echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	client := backendbeer.NewFakeBeerClient(500)
	service := beerapi.NewService(client, s.cfg)
	handler := beerapi.NewHandler(service)

	beers := s.Echo.Group("/beer")
	BeerRoutes(beers, handler)
}
