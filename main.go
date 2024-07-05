package main

import (
	"log"

	"github.com/joho/godotenv"
)

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
	// ------------------------- Starting a new server ------------------------- //

	server := NewAPIserver(":3000", store) // Creating a new instance of APIserver class
	server.Run()                           // Calling the Run() Method on the server instance

	// ---------------------------- xxxxxxxxxxxxxx ----------------------------- //

}
