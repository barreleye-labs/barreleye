package restful

import (
	"github.com/barreleye-labs/barreleye/core"
	"github.com/barreleye-labs/barreleye/core/types"
	"github.com/go-kit/log"
)

type ServerConfig struct {
	Logger     log.Logger
	ListenAddr string
}

type Server struct {
	ServerConfig
	txChan      chan *types.Transaction
	bc          *core.Blockchain
	privateKey  *types.PrivateKey
	faucetLimit map[string]int64 // ip => unix time.
}

func NewServer(cfg ServerConfig, bc *core.Blockchain, txChan chan *types.Transaction, privateKey *types.PrivateKey) *Server {
	return &Server{
		ServerConfig: cfg,
		bc:           bc,
		txChan:       txChan,
		privateKey:   privateKey,
		faucetLimit:  make(map[string]int64),
	}
}
