package restful

import "github.com/labstack/echo/v4"

func (s *Server) Start() error {
	e := echo.New()

	e.GET("/blocks/:id", s.getBlock)
	e.GET("/blocks", s.getBlocks)
	e.GET("/last-block", s.getLastBlock)
	e.GET("txs/:id", s.getTx)
	e.GET("txs", s.getTxs)
	e.GET("/accounts/:address", s.getAccount)
	e.POST("/txs", s.postTx)
	e.POST("/faucet", s.requestSomeCoin)

	return e.Start(s.ListenAddr)
}
