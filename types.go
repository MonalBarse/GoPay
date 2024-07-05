package main

import (
	"time"

	"golang.org/x/exp/rand"
)

// ---------------------------- Defining Account class ---------------------------------- //
type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

// ------------------------------  New Account Creation -------------------------------- //
func NewAccount(fistName, lastName string) *Account {
	return &Account{
		// When NewAccount is called it will return a new instance of Account with the following fields
		FirstName: fistName,
		LastName:  lastName,
		Number:    rand.Int63n(10000000),
		// Balance will automatically be set to 0 if not provided
		CreatedAt: time.Now().UTC(),
	}
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
