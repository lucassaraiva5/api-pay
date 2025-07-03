package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
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
	Amount              int64  `json:"amount"`
	Currency            string `json:"currency"`
	StatementDescriptor string `json:"statementDescriptor"`
	PaymentType         string `json:"paymentType"`
	Card                Card   `json:"card"`
}

type TransactionResponse struct {
	ID                  string `json:"id"`
	Date                string `json:"date"`
	Status              string `json:"status"`
	Amount              int64  `json:"amount"`
	OriginalAmount      int64  `json:"originalAmount"`
	Currency            string `json:"currency"`
	StatementDescriptor string `json:"statementDescriptor"`
	PaymentType         string `json:"paymentType"`
	CardId              string `json:"cardId"`
}

type VoidRequest struct {
	// amount removido, não é mais necessário
}

type transactionInternal struct {
	ID                  string
	Date                string
	Status              string
	Amount              *money.Money
	OriginalAmount      *money.Money
	Currency            string
	StatementDescriptor string
	PaymentType         string
	CardId              string
}

var (
	transactions   = make(map[string]*transactionInternal)
	transactionsMu sync.Mutex
)

func generateID() string {
	return uuid.New().String()
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := generateID()
	cardId := generateID()
	amount := money.New(req.Amount, req.Currency)
	transaction := &transactionInternal{
		ID:                  id,
		Date:                time.Now().Format("2006-01-02"),
		Status:              "paid",
		Amount:              amount,
		OriginalAmount:      amount,
		Currency:            req.Currency,
		StatementDescriptor: req.StatementDescriptor,
		PaymentType:         req.PaymentType,
		CardId:              cardId,
	}
	transactionsMu.Lock()
	transactions[id] = transaction
	transactionsMu.Unlock()
	resp := TransactionResponse{
		ID:                  transaction.ID,
		Date:                transaction.Date,
		Status:              transaction.Status,
		Amount:              transaction.Amount.Amount(),
		OriginalAmount:      transaction.OriginalAmount.Amount(),
		Currency:            transaction.Currency,
		StatementDescriptor: transaction.StatementDescriptor,
		PaymentType:         transaction.PaymentType,
		CardId:              transaction.CardId,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func voidTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req VoidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
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
	if transaction.Status == "voided" {
		http.Error(w, "transaction already voided", http.StatusBadRequest)
		return
	}
	transaction.Status = "voided"
	transaction.Amount = money.New(0, transaction.Currency)
	resp := TransactionResponse{
		ID:                  transaction.ID,
		Date:                transaction.Date,
		Status:              transaction.Status,
		Amount:              transaction.Amount.Amount(),
		OriginalAmount:      transaction.OriginalAmount.Amount(),
		Currency:            transaction.Currency,
		StatementDescriptor: transaction.StatementDescriptor,
		PaymentType:         transaction.PaymentType,
		CardId:              transaction.CardId,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	resp := TransactionResponse{
		ID:                  transaction.ID,
		Date:                transaction.Date,
		Status:              transaction.Status,
		Amount:              transaction.Amount.Amount(),
		OriginalAmount:      transaction.OriginalAmount.Amount(),
		Currency:            transaction.Currency,
		StatementDescriptor: transaction.StatementDescriptor,
		PaymentType:         transaction.PaymentType,
		CardId:              transaction.CardId,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/transactions", createTransactionHandler).Methods("POST")
	r.HandleFunc("/void/{id}", voidTransactionHandler).Methods("POST")
	r.HandleFunc("/transactions/{id}", getTransactionHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	log.Println("Stripe mock server running on :8082")
	http.ListenAndServe(":8082", r)
}
