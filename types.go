package main

import (
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ---------------------------- Defining Account class ---------------------------------- //
type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           float64   `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}
type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

// ------------------------------  New Account Creation -------------------------------- //
/* func NewAccount(fistName, lastName string) *Account {
	return &Account{
		// When NewAccount is called it will return a new instance of Account with the following fields
		FirstName: fistName,
		LastName:  lastName,
		Number:    rand.Int63n(10000000),
		// Balance will automatically be set to 0 if not provided
		CreatedAt: time.Now().UTC(),
	}
}
*/
var (
	globalAccountCounter int64 = 11111111 // Start account numbers from 11111111
	mu                   sync.Mutex
)

// NewAccount creates a new Account instance with a unique account number starting from 11111111
func NewAccount(firstName, lastName, password string) (*Account, error) {
	mu.Lock()
	defer mu.Unlock()

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	accountNumber := globalAccountCounter
	globalAccountCounter++ // Increment the global account counter for the next account

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            accountNumber,
		EncryptedPassword: string(encryptedPassword),
		Balance:           0.0,
		CreatedAt:         time.Now().UTC(),
	}, nil
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

// handle transefer
type TransferRequest struct {
	ToAccountId int `json:"toAccountId"`
	Amount      int `json:"amount"`
}
