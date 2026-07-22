package model


import (
    "errors"
) 


var (
	ErrUserNotFound = errors.New("User not found")
	ErrInvalidPassword = errors.New("Invalid password")

	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrWalletNotFound = errors.New("Wallet not found")
	ErrNotWalletOwner = errors.New("you do not own this wallet")

	ErrNotUserProfileOwner    = errors.New("you are not the owner of this profile")

	ErrInsufficientBalance = errors.New("insufficient balance")
)