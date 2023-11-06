package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/barreleye-labs/barreleye/core"
	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
)

type APIError struct {
	Error string
}

type Block struct {
	Version		  uint32
	DataHash 	  string
	PrevBlockHash string
	Height		  uint32
	Timestamp	  int64
	Validator 	  string
	Signature	  string
}

type ServerConfig struct {
	Logger	   log.Logger
	ListenAddr string
}

type Server struct {
	ServerConfig
	bc *core.Blockchain
}

func NewServer(cfg ServerConfig, bc *core.Blockchain) *Server {
	return &Server{
		ServerConfig: cfg,
		bc:			  bc,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	e.GET("/block/:hashorid", s.handleGetBlock)

	return e.Start(s.ListenAddr)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	hashOrId := c.Param("hashorid")

	height, err := strconv.Atoi(hashOrId)
	if err == nil {
		block, err := s.bc.GetBlock(uint32(height))
		if err != nil {
			// return c.JSON(http.StatusBadRequest, map[string]any{"error": err})
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})	// 위와 같은 의미. 코드 리팩토링
		}

		jsonBlock := Block{
			Version: 	   block.Header.Version,
			Height: 	   block.Header.Height,
			DataHash: 	   block.Header.DataHash.String(),
			PrevBlockHash: block.PrevBlockHash.String(),
			Timestamp: 	   block.Timestamp,
			Validator: 	   block.Validator.Address().String(),
			Signature: 	   block.Signature.String(),
		}

		return c.JSON(http.StatusOK, jsonBlock)
	}

	fmt.Print("errrrrefef:", err)
	fmt.Println()
	panic("need to implement getBlockByHash!!!!")
}