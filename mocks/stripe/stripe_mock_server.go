package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Card struct {
	Number            string `json:"number"`
	Holder            string `json:"holder"`
	CVV               string `json:"cvv"`
	Expiration        string `json:"expiration"`
	InstallmentNumber int    `json:"installmentNumber"`
}

type CreateTransactionRequest struct {
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	StatementDescriptor string  `json:"statementDescriptor"`
	PaymentType         string  `json:"paymentType"`
	Card                Card    `json:"card"`
}

type TransactionResponse struct {
	ID                  string  `json:"id"`
	Date                string  `json:"date"`
	Status              string  `json:"status"`
	Amount              float64 `json:"amount"`
	OriginalAmount      float64 `json:"originalAmount"`
	Currency            string  `json:"currency"`
	StatementDescriptor string  `json:"statementDescriptor"`
	PaymentType         string  `json:"paymentType"`
	CardId              string  `json:"cardId"`
}

type VoidRequest struct {
	Amount float64 `json:"amount"`
}

var (
	transactions   = make(map[string]*TransactionResponse)
	transactionsMu sync.Mutex
)

func generateID() string {
	return RandString(16)
}

func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := generateID()
	cardId := generateID()
	resp := &TransactionResponse{
		ID:                  id,
		Date:                time.Now().Format("2006-01-02"),
		Status:              "paid",
		Amount:              req.Amount,
		OriginalAmount:      req.Amount,
		Currency:            req.Currency,
		StatementDescriptor: req.StatementDescriptor,
		PaymentType:         req.PaymentType,
		CardId:              cardId,
	}
	transactionsMu.Lock()
	transactions[id] = resp
	transactionsMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func voidTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req VoidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	transactionsMu.Lock()
	defer transactionsMu.Unlock()
	transaction, ok := transactions[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	transaction.Status = "voided"
	transaction.Amount -= req.Amount
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func getTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	transactionsMu.Lock()
	defer transactionsMu.Unlock()
	transaction, ok := transactions[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/transactions", createTransactionHandler).Methods("POST")
	r.HandleFunc("/void/{id}", voidTransactionHandler).Methods("POST")
	r.HandleFunc("/transactions/{id}", getTransactionHandler).Methods("GET")
	log.Println("Stripe mock server running on :8082")
	http.ListenAndServe(":8082", r)
}
