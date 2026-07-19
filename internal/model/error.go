package model


import (
    "errors"
    //"github.com/go-playground/validator/v10"
) 


var (
	// 404
	ErrUserNotFound = errors.New("User not found")
	ErrInvalidPassword = errors.New("Invalid password")

	// ErrWalletNotFound = errors.New("Wallet not found")
	// ErrTransNotFound = errors.New("Transaction not found")
)