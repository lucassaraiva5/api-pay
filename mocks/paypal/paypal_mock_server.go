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
	Number         string `json:"number"`
	HolderName     string `json:"holderName"`
	CVV            string `json:"cvv"`
	ExpirationDate string `json:"expirationDate"`
	Installments   int    `json:"installments"`
}

type PaymentMethod struct {
	Type string `json:"type"`
	Card Card   `json:"card"`
}

type CreateChargeRequest struct {
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Description   string        `json:"description"`
	PaymentMethod PaymentMethod `json:"paymentMethod"`
}

type ChargeResponse struct {
	ID             string  `json:"id"`
	CreatedAt      string  `json:"createdAt"`
	Status         string  `json:"status"`
	OriginalAmount float64 `json:"originalAmount"`
	CurrentAmount  float64 `json:"currentAmount"`
	Currency       string  `json:"currency"`
	Description    string  `json:"description"`
	PaymentMethod  string  `json:"paymentMethod"`
	CardId         string  `json:"cardId"`
}

type RefundRequest struct {
	Amount float64 `json:"amount"`
}

var (
	charges   = make(map[string]*ChargeResponse)
	chargesMu sync.Mutex
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

func createChargeHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := generateID()
	cardId := generateID()
	resp := &ChargeResponse{
		ID:             id,
		CreatedAt:      time.Now().Format("2006-01-02"),
		Status:         "authorized",
		OriginalAmount: req.Amount,
		CurrentAmount:  req.Amount,
		Currency:       req.Currency,
		Description:    req.Description,
		PaymentMethod:  req.PaymentMethod.Type,
		CardId:         cardId,
	}
	chargesMu.Lock()
	charges[id] = resp
	chargesMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func refundChargeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req RefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	chargesMu.Lock()
	defer chargesMu.Unlock()
	charge, ok := charges[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	charge.Status = "refunded"
	charge.CurrentAmount -= req.Amount
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charge)
}

func getChargeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	chargesMu.Lock()
	defer chargesMu.Unlock()
	charge, ok := charges[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charge)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/charges", createChargeHandler).Methods("POST")
	r.HandleFunc("/refund/{id}", refundChargeHandler).Methods("POST")
	r.HandleFunc("/charges/{id}", getChargeHandler).Methods("GET")
	log.Println("PayPal mock server running on :8081")
	http.ListenAndServe(":8081", r)
}
