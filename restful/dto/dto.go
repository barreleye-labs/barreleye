package dto

type Signer struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
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
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
}
