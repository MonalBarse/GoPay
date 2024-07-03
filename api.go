package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIserver struct {
	listenAddress string
}

// The apiFunc is basically a function signature of the func we are using to handle the http requests
type apiFunc func(http.ResponseWriter, *http.Request) error

type APIerror struct {
	Error string
}

// makeHTTPHandlerFunc() makes our handle functions (handleAccount, handleGetAccount, handleCreateAccount, handleDeleteAccount, handleTransfer) to be used as http.HandlerFunc
func makeHTTPHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if error := fn(w, r); error != nil {
			//Handle the error
			WriteJSON(w, http.StatusInternalServerError, APIerror{Error: error.Error()}) /*.Error() converts -> string */
		}
	}
}

// Creating a funciton to write JSON response
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data) // NewEncoder takes in io.Writer and lucilly out http.ResponseWriter is an io.Writer it implements the right interface

}

func NewAPIserver(listenAddress string) *APIserver {
	// Creating a new instance of APIserver
	return &APIserver{
		listenAddress: listenAddress,
	}
}

func (server *APIserver) Run() {
	router := mux.NewRouter()

	// router.HandleFunc("/account", server.handleAccount) This does not work as out handleAccount is returning an error and we want a http handler

	router.HandleFunc("/account", makeHTTPHandlerFunc(server.handleAccount))

	log.Println("Server is running on port: ", server.listenAddress)
	http.ListenAndServe(server.listenAddress, router)
}

func (server *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// we are returning an error becuase we dont want to handle the error here, we want to handle it in the Run() function

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

func (server *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (server *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	return nil

}

func (server *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (server *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	return nil
}
