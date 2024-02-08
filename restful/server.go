package restful

import (
	"encoding/hex"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/barreleye-labs/barreleye/restful/dto"
	"math/big"
	"net/http"
	"strconv"

	"github.com/barreleye-labs/barreleye/core"
	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
)

type TxResponse struct {
	TxCount uint     `json:"txCount"`
	Hashes  []string `json:"hashes"`
}

type APIError struct {
	Error string
}

type Block struct {
	Hash          string `json:"hash"`
	Version       uint32 `json:"version"`
	DataHash      string `json:"dataHash"`
	PrevBlockHash string `json:"prevBlockHash"`
	Height        uint32 `json:"height"`
	Timestamp     int64  `json:"timestamp"`
	Validator     string `json:"validator"`
	Signature     string `json:"signature"`

	TxResponse TxResponse `json:"txResponse"`
}

type ServerConfig struct {
	Logger     log.Logger
	ListenAddr string
}

type Server struct {
	txChan chan *types.Transaction
	ServerConfig
	bc *core.Blockchain
}

func NewServer(cfg ServerConfig, bc *core.Blockchain, txChan chan *types.Transaction) *Server {
	return &Server{
		ServerConfig: cfg,
		bc:           bc,
		txChan:       txChan,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	//e.GET("/block/:hashorid", s.handleGetBlock)
	e.GET("/blocks/:id", s.handleGetBlock)
	e.GET("/blocks", s.handleGetBlocks)
	e.GET("/last-block", s.handleGetLastBlock)
	e.GET("txs/:id", s.handleGetTx)
	e.GET("txs", s.handleGetTxs)
	//e.GET("/tx/:hash", s.handleGetTx)
	e.POST("/txs", s.handlePostTx)

	return e.Start(s.ListenAddr)
}

func (s *Server) handleGetLastBlock(c echo.Context) error {
	block, err := s.bc.ReadLastBlock()
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, intoJSONBlock(block))
}

func (s *Server) handleGetTxs(c echo.Context) error {
	query := c.QueryParams()

	page, err := strconv.Atoi(query["page"][0])
	if err != nil {
		return fmt.Errorf("failed to convert page type from string to int")
	}

	size, err := strconv.Atoi(query["size"][0])
	if err != nil {
		return fmt.Errorf("failed to convert size type from string to int")
	}
	txs, err := s.bc.ReadTxs(page, size)
	if err != nil {

		return fmt.Errorf("failed to get txs %s", err)
	}

	return c.JSON(http.StatusOK, txs)
}

func (s *Server) handleGetBlocks(c echo.Context) error {
	query := c.QueryParams()

	page, err := strconv.Atoi(query["page"][0])
	if err != nil {
		return fmt.Errorf("failed to convert page type from string to int")
	}

	size, err := strconv.Atoi(query["size"][0])
	if err != nil {
		return fmt.Errorf("failed to convert size type from string to int")
	}

	blocks, err := s.bc.ReadBlocks(page, size)
	if err != nil {
		return fmt.Errorf("failed to get blocks")
	}

	return c.JSON(http.StatusOK, intoJSONBlocks(blocks))
}

//func (s *Server) handleGetBlocksByHash(c echo.Context) error {
//	query := c.QueryParams()
//
//	b, err := hex.DecodeString(query["hash"][0])
//	if err != nil {
//		return fmt.Errorf("failed to decode string")
//	}
//
//	hash := common.HashFromBytes(b)
//
//	size, err := strconv.Atoi(query["size"][0])
//	if err != nil {
//		return fmt.Errorf("failed to convert size type from string to int")
//	}
//	blocks, err := s.bc.GetBlocks(hash, size)
//	if err != nil {
//		return fmt.Errorf("failed to ")
//	}
//
//	return c.JSON(http.StatusOK, intoJSONBlocks(blocks))
//}

func (s *Server) handlePostTx(c echo.Context) error {
	payload := &dto.TransactionRequest{}
	if err := c.Bind(payload); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	signer := crypto.GetPublicKey(payload.SignerX, payload.SignerY)
	signature := crypto.GetSignature(payload.SignatureR, payload.SignatureS)

	nonceBigInt := new(big.Int)
	nonceBigInt.SetString(payload.Nonce, 16)
	nonce := nonceBigInt.Uint64()

	from, err := hex.DecodeString(payload.From)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid property from")
	}

	to, err := hex.DecodeString(payload.To)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid property to")
	}

	valueBigInt := new(big.Int)
	valueBigInt.SetString(payload.Value, 16)
	value := valueBigInt.Uint64()

	data, err := hex.DecodeString(payload.Data)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid property data")
	}

	tx := types.Transaction{
		Nonce:     nonce,
		From:      common.NewAddressFromBytes(from),
		To:        common.NewAddressFromBytes(to),
		Value:     value,
		Data:      data,
		Signer:    signer,
		Signature: &signature,
	}
	tx.Hash = tx.GetHash()

	s.txChan <- &tx

	return nil
}

//func (s *Server) handleGetTx(c echo.Context) error {
//	hash := c.Param("hash")
//
//	b, err := hex.DecodeString(hash)
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
//	}
//
//	tx, err := s.bc.GetTxByHash(common.HashFromBytes(b))
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
//	}
//
//	return c.JSON(http.StatusOK, tx)
//}

func (s *Server) handleGetTx(c echo.Context) error {
	id := c.Param("id")

	number, err := strconv.Atoi(id)
	// If the error is nil we can assume the height of the block is given.
	if err == nil {
		tx, err := s.bc.ReadTxByNumber(uint32(number))
		if err != nil {
			// return c.JSON(http.StatusBadRequest, map[string]any{"error": err})
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()}) // 위와 같은 의미. 코드 리팩토링
		}

		return c.JSON(http.StatusOK, tx)
	}

	// otherwise assume its the hash

	hash, err := hex.DecodeString(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	tx, err := s.bc.ReadTxByHash(common.HashFromBytes(hash))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, tx)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	id := c.Param("id")

	height, err := strconv.Atoi(id)
	// If the error is nil we can assume the height of the block is given.
	if err == nil {
		block, err := s.bc.ReadBlockByHeight(uint32(height))
		if err != nil {
			// return c.JSON(http.StatusBadRequest, map[string]any{"error": err})
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()}) // 위와 같은 의미. 코드 리팩토링
		}

		return c.JSON(http.StatusOK, intoJSONBlock(block))
	}

	// otherwise assume its the hash

	b, err := hex.DecodeString(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	block, err := s.bc.ReadBlockByHash(common.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, intoJSONBlock(block))
}

func intoJSONBlock(block *types.Block) Block {
	txResponse := TxResponse{
		TxCount: uint(len(block.Transactions)),
		Hashes:  make([]string, len(block.Transactions)),
	}

	for i := 0; i < int(txResponse.TxCount); i++ {
		txResponse.Hashes[i] = block.Transactions[i].GetHash().String()
	}

	return Block{
		Hash:          block.GetHash().String(),
		Version:       block.Header.Version,
		Height:        block.Header.Height,
		DataHash:      block.Header.DataHash.String(),
		PrevBlockHash: block.PrevBlockHash.String(),
		Timestamp:     block.Timestamp,
		Validator:     block.Validator.Address().String(),
		Signature:     block.Signature.String(),
		TxResponse:    txResponse,
	}
}

func intoJSONBlocks(blocks []*types.Block) []Block {
	value := []Block{}
	for i := 0; i < len(blocks); i++ {
		txResponse := TxResponse{
			TxCount: uint(len(blocks[i].Transactions)),
			Hashes:  make([]string, len(blocks[i].Transactions)),
		}

		for j := 0; j < int(txResponse.TxCount); j++ {
			txResponse.Hashes[j] = blocks[i].Transactions[j].GetHash().String()
		}

		b := Block{
			Hash:          blocks[i].GetHash().String(),
			Version:       blocks[i].Header.Version,
			Height:        blocks[i].Header.Height,
			DataHash:      blocks[i].Header.DataHash.String(),
			PrevBlockHash: blocks[i].PrevBlockHash.String(),
			Timestamp:     blocks[i].Timestamp,
			Validator:     blocks[i].Validator.Address().String(),
			Signature:     blocks[i].Signature.String(),
			TxResponse:    txResponse,
		}

		value = append(value, b)
	}
	return value
}
