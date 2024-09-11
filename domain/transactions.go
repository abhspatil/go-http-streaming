package domain

import "time"

type TransactionHistoryData struct {
	TransactionId string    `json:"transaction_id"`
	Amount        string    `json:"amount"`
	Network       string    `json:"network"`
	CreatedOn     time.Time `json:"created_on"`
}
