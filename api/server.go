package api

import (
	"net/http"

	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
)

type ServerConfig struct {
	Logger	   log.Logger
	ListenAddr string
}

type Server struct {
	ServerConfig
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
	}
}

func (s *Server) start() error {
	e := echo.New()

	e.GET("/block/:hashorid", s.handleGetBlock)

	return e.Start(s.ListenAddr)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"msg": "it works!"})
}