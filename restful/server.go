package restful

import (
	"encoding/hex"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/common/util"
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
	Height        int32  `json:"height"`
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
	e.GET("/blocks/:id", s.getBlock)
	e.GET("/blocks", s.getBlocks)
	e.GET("/last-block", s.getLastBlock)
	e.GET("txs/:id", s.getTx)
	e.GET("txs", s.getTxs)
	e.GET("/accounts/:address", s.getAccount)
	//e.GET("/tx/:hash", s.handleGetTx)
	e.POST("/txs", s.postTx)

	return e.Start(s.ListenAddr)
}

func (s *Server) getAccount(c echo.Context) error {
	address := c.Param("address")

	var result *types.Account = nil

	bytes, err := hex.DecodeString(address)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	result, err = s.bc.ReadAccountByAddress(common.NewAddressFromBytes(bytes))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.AccountResponse{Account: dto.Account{
		Address: result.Address.String(),
		Balance: hex.EncodeToString(util.Uint64ToBytes(result.Balance)),
	}})
}

func (s *Server) getLastBlock(c echo.Context) error {
	block, err := s.bc.ReadLastBlock()
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, intoJSONBlock(block))
}

func (s *Server) getTxs(c echo.Context) error {
	query := c.QueryParams()

	page, err := strconv.Atoi(query["page"][0])
	if err != nil {
		return fmt.Errorf("failed to convert page type from string to int")
	}

	size, err := strconv.Atoi(query["size"][0])
	if err != nil {
		return fmt.Errorf("failed to convert size type from string to int")
	}

	result, err := s.bc.ReadTxs(page, size)
	if err != nil {

		return fmt.Errorf("failed to get txs %s", err)
	}

	txs := []dto.Transaction{}
	for i := 0; i < len(result); i++ {
		signer := dto.Signer{
			X: hex.EncodeToString(result[i].Signer.Key.X.Bytes()),
			Y: hex.EncodeToString(result[i].Signer.Key.Y.Bytes()),
		}

		signature := dto.Signature{
			R: hex.EncodeToString(result[i].Signature.R.Bytes()),
			S: hex.EncodeToString(result[i].Signature.S.Bytes()),
		}

		tx := dto.Transaction{
			Hash:      result[i].Hash.String(),
			Nonce:     hex.EncodeToString(util.Uint64ToBytes(result[i].Nonce)),
			From:      result[i].From.String(),
			To:        result[i].To.String(),
			Value:     hex.EncodeToString(util.Uint64ToBytes(result[i].Value)),
			Data:      hex.EncodeToString(result[i].Data),
			Signer:    signer,
			Signature: signature,
		}

		txs = append(txs, tx)
	}

	return c.JSON(http.StatusOK, dto.TransactionsResponse{Transactions: txs})
}

func (s *Server) getBlocks(c echo.Context) error {
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

func (s *Server) postTx(c echo.Context) error {
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

	return c.JSON(http.StatusOK, dto.TransactionResponse{Transaction: dto.Transaction{
		Hash:  tx.Hash.String(),
		Nonce: payload.Nonce,
		From:  payload.From,
		To:    payload.To,
		Value: payload.Value,
		Data:  payload.Data,
		Signer: dto.Signer{
			X: payload.SignerX,
			Y: payload.SignerY,
		},
		Signature: dto.Signature{
			R: payload.SignatureR,
			S: payload.SignatureS,
		},
	}})
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

func (s *Server) getTx(c echo.Context) error {
	id := c.Param("id")

	var result *types.Transaction = nil

	number, err := strconv.Atoi(id)
	if err == nil {
		result, err = s.bc.ReadTxByNumber(uint32(number))
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()}) // 위와 같은 의미. 코드 리팩토링
		}
	} else {
		hash, err := hex.DecodeString(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
		}

		result, err = s.bc.ReadTxByHash(common.HashFromBytes(hash))
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}

	signer := dto.Signer{
		X: hex.EncodeToString(result.Signer.Key.X.Bytes()),
		Y: hex.EncodeToString(result.Signer.Key.Y.Bytes()),
	}

	signature := dto.Signature{
		R: hex.EncodeToString(result.Signature.R.Bytes()),
		S: hex.EncodeToString(result.Signature.S.Bytes()),
	}

	tx := dto.Transaction{
		Hash:      result.Hash.String(),
		Nonce:     hex.EncodeToString(util.Uint64ToBytes(result.Nonce)),
		From:      result.From.String(),
		To:        result.To.String(),
		Value:     hex.EncodeToString(util.Uint64ToBytes(result.Value)),
		Data:      hex.EncodeToString(result.Data),
		Signer:    signer,
		Signature: signature,
	}

	return c.JSON(http.StatusOK, dto.TransactionResponse{Transaction: tx})
}

func (s *Server) getBlock(c echo.Context) error {
	id := c.Param("id")

	height, err := strconv.Atoi(id)
	if err == nil {
		block, err := s.bc.ReadBlockByHeight(int32(height))
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, intoJSONBlock(block))
	}

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
