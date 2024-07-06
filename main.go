package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func seedAccount(store Storage, fname, lname, pass string) *Account {
	acc, err := NewAccount(fname, lname, pass)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedData(store Storage) {
	seedAccount(store, "John", "Doe", "password")
	seedAccount(store, "Jane", "Doe", "password")
	seedAccount(store, "Alice", "Bob", "password")
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// ------------------ Database conn and store initialization ---------------- //

	store, err := NewPostgresStore() // Creating a new instance of PostgresStore
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil { // Init() - Initializes the db and creates the table (if not exists)
		log.Fatal(err)
	}

	// ---------------------------- xxxxxxxxxxxxxx ----------------------------- //

	// ---------------------------- Seeding the database ----------------------------- //
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()
	if *seed {
		seedData(store)
	}

	// ------------------------- Starting a new server ------------------------- //

	server := NewAPIserver(":3000", store) // Creating a new instance of APIserver class
	server.Run()                           // Calling the Run() Method on the server instance

	// ---------------------------- xxxxxxxxxxxxxx ----------------------------- //

}
