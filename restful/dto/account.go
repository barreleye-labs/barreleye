package dto

type Account struct {
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	Balance string `json:"balance"`
}

type AccountResponse struct {
	Account Account `json:"account"`
}
