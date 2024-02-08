package dto

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
