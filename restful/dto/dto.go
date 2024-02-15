package dto

type Signer struct {
	X string `json:"x"`
	Y string `json:"y"`
}

func CreateSigner(x string, y string) Signer {
	return Signer{
		X: x,
		Y: y,
	}
}

type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
}

func CreateSignature(r string, s string) Signature {
	return Signature{
		R: r,
		S: s,
	}
}

type Block struct {
	Hash          string    `json:"hash"`
	Version       uint32    `json:"version"`
	DataHash      string    `json:"dataHash"`
	PrevBlockHash string    `json:"prevBlockHash"`
	Height        int32     `json:"height"`
	Timestamp     int64     `json:"timestamp"`
	Signer        string    `json:"signer"`
	Signature     Signature `json:"signature"`
	TxCount       uint32    `json:"txCount"`
	Transactions  []string  `json:"transactions"`
}

type Transaction struct {
	Hash      string    `json:"hash"`
	Nonce     string    `json:"nonce"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Value     string    `json:"value"`
	Data      string    `json:"data"`
	Signer    Signer    `json:"signer"`
	Signature Signature `json:"signature"`
}

func CreateTransaction(
	hash string,
	nonce string,
	from string,
	to string,
	value string,
	data string,
	signer Signer,
	signature Signature) Transaction {
	return Transaction{
		Hash:      hash,
		Nonce:     nonce,
		From:      from,
		To:        to,
		Value:     value,
		Data:      data,
		Signer:    signer,
		Signature: signature,
	}
}

type Account struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type AccountResponse struct {
	Account Account `json:"account"`
}

type TransactionRequest struct {
	Nonce      string `json:"nonce"`
	From       string `json:"from"`
	To         string `json:"to"`
	Value      string `json:"value"`
	Data       string `json:"data"`
	SignerX    string `json:"signerX"`
	SignerY    string `json:"signerY"`
	SignatureR string `json:"signatureR"`
	SignatureS string `json:"signatureS"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	TotalCount   uint32        `json:"totalCount"`
}

func CreateTransactionsResponse(transactions []Transaction, totalCount uint32) TransactionsResponse {
	return TransactionsResponse{
		Transactions: transactions,
		TotalCount:   totalCount,
	}
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
}

func CreateTransactionResponse(transaction Transaction) TransactionResponse {
	return TransactionResponse{
		Transaction: transaction,
	}
}

type BlockResponse struct {
	Block Block `json:"block"`
}

func CreateBlockResponse(block Block) BlockResponse {
	return BlockResponse{
		Block: block,
	}
}

type BlocksResponse struct {
	Blocks     []Block `json:"blocks"`
	TotalCount uint32  `json:"totalCount"`
}

func CreateBlocksResponse(blocks []Block, totalCount uint32) BlocksResponse {
	return BlocksResponse{
		Blocks:     blocks,
		TotalCount: totalCount,
	}
}

func CreateBlock(
	hash string,
	version uint32,
	dataHash string,
	prevBlockHash string,
	height int32,
	timestamp int64,
	signer string,
	signature Signature,
	txCount uint32,
	transactions []string) Block {
	return Block{
		Hash:          hash,
		Version:       version,
		DataHash:      dataHash,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     timestamp,
		Signer:        signer,
		Signature:     signature,
		TxCount:       txCount,
		Transactions:  transactions,
	}
}
