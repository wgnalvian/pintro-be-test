package dto

type TransferRequest struct {
	Amount int    `json:"amount"`
	To     string `json:"to"`
}
