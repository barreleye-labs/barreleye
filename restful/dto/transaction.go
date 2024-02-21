package dto

type Transaction struct {
	Hash        string    `json:"hash"`
	Nonce       string    `json:"nonce"`
	BlockHeight int32     `json:"blockHeight"`
	Timestamp   int64     `json:"timestamp"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       string    `json:"value"`
	Data        string    `json:"data"`
	Signer      Signer    `json:"signer"`
	Signature   Signature `json:"signature"`
}

func CreateTransaction(
	hash string,
	nonce string,
	blockHeight int32,
	timestamp int64,
	from string,
	to string,
	value string,
	data string,
	signer Signer,
	signature Signature) Transaction {
	return Transaction{
		Hash:        hash,
		Nonce:       nonce,
		BlockHeight: blockHeight,
		Timestamp:   timestamp,
		From:        from,
		To:          to,
		Value:       value,
		Data:        data,
		Signer:      signer,
		Signature:   signature,
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

type FaucetRequest struct {
	AccountAddress string `json:"accountAddress"`
}
