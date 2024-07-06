package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
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

	router.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin))

	router.HandleFunc("/account", makeHTTPHandlerFunc(server.handleAccount))

	router.HandleFunc("/account/{id}", requireJWT(makeHTTPHandlerFunc(server.handleGetAccountById), server.store))
	router.HandleFunc("/transfer", makeHTTPHandlerFunc(server.handleTransfer))

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
  4. getId()                    -> to get the id from the request
*/

type APIerror struct {
	Error string `json:"error"`
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

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func requireJWT(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			WriteJSON(w, http.StatusUnauthorized, APIerror{Error: "Missing token"})
			return
		}

		token, err := verifyJWT(tokenString)
		if err != nil {
			log.Printf("JWT verification error: %v", err)
			WriteJSON(w, http.StatusUnauthorized, APIerror{Error: "Unauthorized"})
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, APIerror{Error: "Invalid token"})
			return
		}

		userID, err := getID(r)
		if err != nil {
			log.Printf("Error extracting user ID: %v", err)
			WriteJSON(w, http.StatusBadRequest, APIerror{Error: "Invalid user ID"})
			return
		}

		account, err := s.GetAccountbyId(userID)
		if err != nil {
			log.Printf("Error fetching account: %v", err)
			WriteJSON(w, http.StatusNotFound, APIerror{Error: "Account not found"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			WriteJSON(w, http.StatusUnauthorized, APIerror{Error: "Unauthorized"})
			return
		}

		// Ensure the requester owns the account
		expectedAccountNumber, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["accountNumber"]), 10, 64)
		if err != nil || account.Number != expectedAccountNumber {
			WriteJSON(w, http.StatusForbidden, APIerror{Error: "Forbidden"})
			return
		}

		// Add accountNumber to request context
		ctx := context.WithValue(r.Context(), "accountNumber", account.Number)
		handlerFunc(w, r.WithContext(ctx))
	}
}

func verifyJWT(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func createJWT(account *Account) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"expiresAt":     time.Now().Add(time.Hour * 24).Unix(),
		"accountNumber": account.Number,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// ---------------------------------------- xxxxxxxxxx ------------------------------------------------ //

// ------------------------------- Methods of APIserver class ---------------------------------------- //
/*
Index :
  1. handleGetAccount()         -> to handle Get account
  2. handleGetAccountById()     -> to handle Get account by id
  3. handleCreateAccount()      -> to handle Create account
  4. handleDeleteAccount()      -> to handle Delete account
  5. handleTransfer()           -> to handle Transfer
  6. handleLogin()              -> to handle Login
*/

func (server *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := server.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

// 2
func (server *APIserver) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	// Retrieve accountNumber from context
	accountNumber, ok := r.Context().Value("accountNumber").(int64)
	if !ok {
		return fmt.Errorf("failed to get accountNumber from context")
	}

	switch r.Method {
	case "GET":
		account, err := server.store.GetAccountbyId(id)
		if err != nil {
			return err
		}

		// Ensure the requester owns the account
		if account.Number != accountNumber {
			return WriteJSON(w, http.StatusForbidden, APIerror{Error: "Forbidden"})
		}

		return WriteJSON(w, http.StatusOK, account)

	case "DELETE":
		account, err := server.store.GetAccountbyId(id)
		if err != nil {
			return err
		}

		// Ensure the requester owns the account
		if account.Number != accountNumber {
			return WriteJSON(w, http.StatusForbidden, APIerror{Error: "Forbidden"})
		}

		if err := server.store.DeleteAccount(id); err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})

	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

// 2
func (server *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}
	account, err := NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName, createAccountRequest.Password)
	if err != nil {
		return err
	}
	if err := server.store.CreateAccount(account); err != nil {
		return err
	}
	tokenStr, err := createJWT(account)
	if err != nil {
		return err
	}
	fmt.Println("Token: ", tokenStr)
	return WriteJSON(w, http.StatusOK, tokenStr)
}

// 4
func (server *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := server.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})

}

// 5
func (server *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferRequest)
}

// 6
func (server *APIserver) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	acc, err := server.store.GetAccountByNumber(req.Number)
	if err != nil {
		return err
	}

	fmt.Println("Account: ", acc)

	return WriteJSON(w, http.StatusOK, req)
}

// ---------------------------------------- xxxxxxxxxx ------------------------------------------------ //
