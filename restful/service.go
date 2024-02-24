package restful

import (
	"encoding/hex"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/common/util"
	"github.com/barreleye-labs/barreleye/config"
	"github.com/barreleye-labs/barreleye/core/types"
	"github.com/barreleye-labs/barreleye/restful/dto"
	"github.com/labstack/echo/v4"
	"math/big"
	"net/http"
	"strconv"
)

func (s *Server) requestSomeCoin(c echo.Context) error {
	payload := &dto.FaucetRequest{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid payload "+err.Error()))
	}

	to, err := hex.DecodeString(payload.AccountAddress)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid AccountAddress "+err.Error())
	}

	accountNonce, err := s.bc.ReadAccountNonceByAddress(s.privateKey.PublicKey.Address())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
	}

	nonce := uint64(0)
	if accountNonce != nil {
		nonce = *accountNonce
	}

	tx := types.CreateTransaction(
		nonce,
		s.privateKey.PublicKey.Address(),
		common.NewAddressFromBytes(to),
		config.FaucetAmount,
		[]byte{171})

	tx.Hash = tx.GetHash()
	if err = tx.Sign(s.privateKey); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
	}

	s.txChan <- tx

	signerDTO := dto.CreateSigner(tx.Signer.Key.X.Text(16), tx.Signer.Key.Y.Text(16))
	signatureDTO := dto.CreateSignature(tx.Signature.R.Text(16), tx.Signature.S.Text(16))

	txDTO := dto.CreateTransaction(
		tx.Hash.String(),
		hex.EncodeToString(util.Uint64ToBytes(tx.Nonce)),
		-1,
		-1,
		tx.From.String(),
		tx.To.String(),
		hex.EncodeToString(util.Uint64ToBytes(tx.Value)),
		hex.EncodeToString(tx.Data),
		signerDTO,
		signatureDTO)

	return c.JSON(http.StatusOK, ResponseOk(dto.CreateTransactionResponse(txDTO)))
}

func (s *Server) getAccount(c echo.Context) error {
	address := c.Param("address")

	var result *types.Account = nil

	bytes, err := hex.DecodeString(address)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid address"))
	}

	result, err = s.bc.ReadAccountByAddress(common.NewAddressFromBytes(bytes))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
	}

	if result == nil {
		return c.JSON(http.StatusNotFound, ResponseNotFound("not found account"))
	}

	return c.JSON(http.StatusOK, ResponseOk(dto.AccountResponse{Account: dto.Account{
		Address: result.Address.String(),
		Nonce:   hex.EncodeToString(util.Uint64ToBytes(result.Nonce)),
		Balance: hex.EncodeToString(util.Uint64ToBytes(result.Balance)),
	}}))
}

func (s *Server) getLastBlock(c echo.Context) error {
	result, err := s.bc.ReadLastBlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
	}

	if result == nil {
		return c.JSON(http.StatusNotFound, ResponseNotFound("not found last block"))
	}

	transactions := []string{}
	for j := 0; j < len(result.Transactions); j++ {
		transactions = append(transactions, result.Transactions[j].Hash.String())
	}

	signature := dto.CreateSignature(result.Signature.R.Text(16), result.Signature.S.Text(16))

	block := dto.CreateBlock(
		result.Hash.String(),
		result.Version,
		result.DataHash.String(),
		result.PrevBlockHash.String(),
		result.Height,
		result.Timestamp,
		result.Signer.Address().String(),
		result.Extra,
		signature,
		uint32(len(result.Transactions)),
		transactions)
	return c.JSON(http.StatusOK, ResponseOk(dto.CreateBlockResponse(block)))
}

func (s *Server) getTxs(c echo.Context) error {
	query := c.QueryParams()

	page, err := strconv.Atoi(query["page"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid page"))
	}

	size, err := strconv.Atoi(query["size"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid size"))
	}

	result, err := s.bc.ReadTxs(page, size)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
	}

	txs := []dto.Transaction{}
	for i := 0; i < len(result); i++ {
		signer := dto.CreateSigner(result[i].Signer.Key.X.Text(16), result[i].Signer.Key.Y.Text(16))
		signature := dto.CreateSignature(result[i].Signature.R.Text(16), result[i].Signature.S.Text(16))

		tx := dto.CreateTransaction(
			result[i].Hash.String(),
			hex.EncodeToString(util.Uint64ToBytes(result[i].Nonce)),
			result[i].BlockHeight,
			result[i].Timestamp,
			result[i].From.String(),
			result[i].To.String(),
			hex.EncodeToString(util.Uint64ToBytes(result[i].Value)),
			hex.EncodeToString(result[i].Data),
			signer,
			signature)

		txs = append(txs, tx)
	}

	lastTxNumber, err := s.bc.ReadLastTxNumber()
	if err != nil {
		return err
	}

	totalCount := uint32(0)
	if lastTxNumber != nil {
		totalCount = *lastTxNumber + 1
	}

	return c.JSON(http.StatusOK, ResponseOk(dto.CreateTransactionsResponse(txs, totalCount)))
}

func (s *Server) getBlocks(c echo.Context) error {
	query := c.QueryParams()

	page, err := strconv.Atoi(query["page"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid page"))
	}

	size, err := strconv.Atoi(query["size"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid size"))
	}

	result, err := s.bc.ReadBlocks(page, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
	}

	blocks := []dto.Block{}
	for i := 0; i < len(result); i++ {
		transactions := []string{}
		for j := 0; j < len(result[i].Transactions); j++ {
			transactions = append(transactions, result[i].Transactions[j].Hash.String())
		}

		signature := dto.CreateSignature(result[i].Signature.R.Text(16), result[i].Signature.S.Text(16))

		block := dto.CreateBlock(
			result[i].Hash.String(),
			result[i].Version,
			result[i].DataHash.String(),
			result[i].PrevBlockHash.String(),
			result[i].Height,
			result[i].Timestamp,
			result[i].Signer.Address().String(),
			result[i].Extra,
			signature,
			uint32(len(result[i].Transactions)),
			transactions)

		blocks = append(blocks, block)
	}

	lastBlockHeight, err := s.bc.ReadLastBlockHeight()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ResponseOk(dto.CreateBlocksResponse(blocks, uint32(*lastBlockHeight+1))))
}

func (s *Server) postTx(c echo.Context) error {
	payload := &dto.TransactionRequest{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid payload "+err.Error()))
	}

	signer := types.GetPublicKey(payload.SignerX, payload.SignerY)
	signature := types.GetSignature(payload.SignatureR, payload.SignatureS)

	from, err := hex.DecodeString(payload.From)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid from "+err.Error()))
	}

	to, err := hex.DecodeString(payload.To)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid to "+err.Error()))
	}

	accountNonce, err := s.bc.ReadAccountNonceByAddress(common.NewAddressFromBytes(from))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
	}

	nonce := uint64(0)
	if accountNonce != nil {
		nonce = *accountNonce
	}

	valueBigInt := new(big.Int)
	valueBigInt.SetString(payload.Value, 16)
	value := valueBigInt.Uint64()

	data, err := hex.DecodeString(payload.Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid data "+err.Error()))
	}

	tx := types.CreateSignedTransaction(
		nonce,
		common.NewAddressFromBytes(from),
		common.NewAddressFromBytes(to),
		value,
		data,
		*signer,
		signature)

	tx.Hash = tx.GetHash()

	s.txChan <- tx

	signerDTO := dto.CreateSigner(payload.SignerX, payload.SignerY)
	signatureDTO := dto.CreateSignature(payload.SignatureR, payload.SignatureS)

	txDTO := dto.CreateTransaction(
		tx.Hash.String(),
		hex.EncodeToString(util.Uint64ToBytes(tx.Nonce)),
		-1,
		-1,
		payload.From,
		payload.To,
		payload.Value,
		payload.Data,
		signerDTO,
		signatureDTO)

	return c.JSON(http.StatusOK, ResponseOk(dto.CreateTransactionResponse(txDTO)))
}

func (s *Server) getTx(c echo.Context) error {
	id := c.Param("id")

	var result *types.Transaction = nil

	number, err := strconv.Atoi(id)
	if err == nil {
		result, err = s.bc.ReadTxByNumber(uint32(number))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
		}
	} else {
		hash, err := hex.DecodeString(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid hash "+id))
		}

		result, err = s.bc.ReadTxByHash(common.HashFromBytes(hash))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
		}
	}

	if result == nil {
		return c.JSON(http.StatusBadRequest, ResponseNotFound("not found transaction"+id))
	}

	signer := dto.CreateSigner(result.Signer.Key.X.Text(16), result.Signer.Key.Y.Text(16))
	signature := dto.CreateSignature(result.Signature.R.Text(16), result.Signature.S.Text(16))

	tx := dto.CreateTransaction(
		result.Hash.String(),
		hex.EncodeToString(util.Uint64ToBytes(result.Nonce)),
		result.BlockHeight,
		result.Timestamp,
		result.From.String(),
		result.To.String(),
		hex.EncodeToString(util.Uint64ToBytes(result.Value)),
		hex.EncodeToString(result.Data),
		signer,
		signature)

	return c.JSON(http.StatusOK, ResponseOk(dto.CreateTransactionResponse(tx)))
}

func (s *Server) getBlock(c echo.Context) error {
	id := c.Param("id")

	var result *types.Block = nil

	height, err := strconv.Atoi(id)
	if err == nil {
		result, err = s.bc.ReadBlockByHeight(int32(height))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
		}
	} else {
		hash, err := hex.DecodeString(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseBadRequest("invalid hash "+id))
		}

		result, err = s.bc.ReadBlockByHash(common.HashFromBytes(hash))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseServerError(err.Error()))
		}
	}

	if result == nil {
		return c.JSON(http.StatusNotFound, ResponseNotFound("not found block"))
	}

	transactions := []string{}
	for j := 0; j < len(result.Transactions); j++ {
		transactions = append(transactions, result.Transactions[j].Hash.String())
	}

	signature := dto.CreateSignature(result.Signature.R.Text(16), result.Signature.S.Text(16))

	block := dto.CreateBlock(
		result.Hash.String(),
		result.Version,
		result.DataHash.String(),
		result.PrevBlockHash.String(),
		result.Height,
		result.Timestamp,
		result.Signer.Address().String(),
		result.Extra,
		signature,
		uint32(len(result.Transactions)),
		transactions)
	return c.JSON(http.StatusOK, ResponseOk(dto.CreateBlockResponse(block)))
}
