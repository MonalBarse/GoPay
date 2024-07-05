package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// ----------------------------------- Defining APIserver class ----------------------------------- //
/*
Index : APIserver is a struct with - listenAddress and two methods - Run() and handleAccount()
  1. APIserver struct           -> is just listening address for the server and a store instance
  2. Run() function             -> to run the server
  3. handleAccount() function   -> to handle all the account related requests
*/

type APIserver struct {
	listenAddress string
	store         Storage
}

func (server *APIserver) Run() {
	router := mux.NewRouter()

	// router.HandleFunc("/account", server.handleAccount) This does not work as out handleAccount is returning an error and we want a http handler

	router.HandleFunc("/account", makeHTTPHandlerFunc(server.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(server.handleGetAccountById))

	log.Println("Server is running on port: ", server.listenAddress)
	http.ListenAndServe(server.listenAddress, router)
}

func (server *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {

	var method = r.Method
	switch method {
	case "GET":
		return server.handleGetAccount(w, r)
	case "POST":
		return server.handleCreateAccount(w, r)
	case "DELETE":
		return server.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed %s", method)
	}

}

// ---------------------------------------- xxxxxxxxxx ------------------------------------------------ //

// ------------------------------------- Creating API server ------------------------------------------ //

// Creating a new instance of APIserver - NewAPIserver
func NewAPIserver(listenAddress string, store Storage) *APIserver {
	// Creating a new instance of APIserver
	return &APIserver{
		listenAddress: listenAddress,
		store:         store,
	}
}

// ---------------------------------------- xxxxxxxxxx ------------------------------------------------ //

// ------------------------------------- Utility Functions -------------------------------------------- //

/*
Index :
  1. APIerror struct            -> to handle the error response of the
  2. apiFunc type               -> to define the function signature of the functions we are using to handle the http requests
  4. WriteJSON()                -> to write JSON response
  3. makeHTTPHandlerFunc()      -> to make our handle functions (handleAccount, handleGetAccount,....) to be used as http.HandlerFunc
*/

type APIerror struct {
	Error string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data) // NewEncoder takes in io.Writer and lucilly out http.ResponseWriter is an io.Writer
}

func makeHTTPHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if error := fn(w, r); error != nil {
			//Handle the error
			WriteJSON(w, http.StatusInternalServerError, APIerror{Error: error.Error()}) /*.Error() converts -> string */
		}
	}
}

// ---------------------------------------- xxxxxxxxxx ------------------------------------------------ //

// ------------------------------ Handler Functions (type apiFunc) ------------------------------------ //

/*
Index :
  1. handleGetAccount()         -> to handle Get account
  2. handleCreateAccount()      -> to handle Create account
  3. handleDeleteAccount()      -> to handle Delete account
  4. handleTransfer()           -> to handle Transfer
*/

func (server *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := server.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

// 1
func (server *APIserver) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	id := mux.Vars(r)["id"]
	fmt.Println(id)
	return WriteJSON(w, http.StatusOK, &Account{})
}

// 2
func (server *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	createAccountRequest := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}
	account := NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	if err := server.store.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

// 3
func (server *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// 4
func (server *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// ---------------------------------------- xxxxxxxxxx ------------------------------------------------ //
