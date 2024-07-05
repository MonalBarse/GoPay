package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// -------------------------- Storage Interface --------------------------- //
type Storage interface { // Storage interface - an instance of this interface will be responsible for all the database operations
	CreateAccount(*Account) error // CreateAccount() - will takes in an account and persist it to the database
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountbyId(int) (*Account, error)
}

// ------------ PostgresStore Class - persists data in DB --------------- //

type PostgresStore struct {
	db *sql.DB // db - used to execute SQL queries against the Postgres database
}

// Methods for PostgresStore
func (s *PostgresStore) createAccountTable() error {
	query := `
      create table if not exists accounts(
      id serial primary key,
      first_name varchar(50),
      last_name varchar(50),
      number bigint unique not null,  
      balance FLOAT,
      created_at timestamp default current_timestamp)
    `

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
        INSERT INTO accounts (first_name, last_name, number, balance, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	res, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", res)
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		if err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}
	return accounts, nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) GetAccountbyId(id int) (*Account, error) {
	return nil, nil
}

// ---------------------------- xxxxxxxxxxxxxx ----------------------------- //

// ---------------------- PostgresStore Constructor ------------------------ //
func NewPostgresStore() (*PostgresStore, error) {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	// Construct connection string
	connectionString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s",
		dbUser, dbName, dbPassword, dbSSLMode)

	// Open database connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Ping database to check connectivity
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

// ---------------------------- xxxxxxxxxxxxxx ----------------------------- //
