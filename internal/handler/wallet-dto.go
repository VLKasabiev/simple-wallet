package handler

type CreateWalletRequest struct {
	UserID   int    `json:"user_id"`
	Currency string `json:"currency"`
}