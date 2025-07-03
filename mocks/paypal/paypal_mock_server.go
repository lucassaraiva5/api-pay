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
	Amount        int64         `json:"amount"`
	Currency      string        `json:"currency"`
	Description   string        `json:"description"`
	PaymentMethod PaymentMethod `json:"paymentMethod"`
}

type ChargeResponse struct {
	ID             string `json:"id"`
	CreatedAt      string `json:"createdAt"`
	Status         string `json:"status"`
	OriginalAmount int64  `json:"originalAmount"`
	CurrentAmount  int64  `json:"currentAmount"`
	Currency       string `json:"currency"`
	Description    string `json:"description"`
	PaymentMethod  string `json:"paymentMethod"`
	CardId         string `json:"cardId"`
}

type RefundRequest struct {
	// amount removido, não é mais necessário
}

type chargeInternal struct {
	ID             string
	CreatedAt      string
	Status         string
	OriginalAmount *money.Money
	CurrentAmount  *money.Money
	Currency       string
	Description    string
	PaymentMethod  string
	CardId         string
}

var (
	charges   = make(map[string]*chargeInternal)
	chargesMu sync.Mutex
)

func generateID() string {
	return uuid.New().String()
}

func createChargeHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := generateID()
	cardId := generateID()
	amount := money.New(req.Amount, req.Currency)
	charge := &chargeInternal{
		ID:             id,
		CreatedAt:      time.Now().Format("2006-01-02"),
		Status:         "authorized",
		OriginalAmount: amount,
		CurrentAmount:  amount,
		Currency:       req.Currency,
		Description:    req.Description,
		PaymentMethod:  req.PaymentMethod.Type,
		CardId:         cardId,
	}
	chargesMu.Lock()
	charges[id] = charge
	chargesMu.Unlock()
	resp := ChargeResponse{
		ID:             charge.ID,
		CreatedAt:      charge.CreatedAt,
		Status:         charge.Status,
		OriginalAmount: charge.OriginalAmount.Amount(),
		CurrentAmount:  charge.CurrentAmount.Amount(),
		Currency:       charge.Currency,
		Description:    charge.Description,
		PaymentMethod:  charge.PaymentMethod,
		CardId:         charge.CardId,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func refundChargeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req RefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
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
	if charge.Status == "refunded" {
		http.Error(w, "charge already refunded", http.StatusBadRequest)
		return
	}
	charge.Status = "refunded"
	charge.CurrentAmount = money.New(0, charge.Currency)
	resp := ChargeResponse{
		ID:             charge.ID,
		CreatedAt:      charge.CreatedAt,
		Status:         charge.Status,
		OriginalAmount: charge.OriginalAmount.Amount(),
		CurrentAmount:  charge.CurrentAmount.Amount(),
		Currency:       charge.Currency,
		Description:    charge.Description,
		PaymentMethod:  charge.PaymentMethod,
		CardId:         charge.CardId,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	resp := ChargeResponse{
		ID:             charge.ID,
		CreatedAt:      charge.CreatedAt,
		Status:         charge.Status,
		OriginalAmount: charge.OriginalAmount.Amount(),
		CurrentAmount:  charge.CurrentAmount.Amount(),
		Currency:       charge.Currency,
		Description:    charge.Description,
		PaymentMethod:  charge.PaymentMethod,
		CardId:         charge.CardId,
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
	r.HandleFunc("/charges", createChargeHandler).Methods("POST")
	r.HandleFunc("/refund/{id}", refundChargeHandler).Methods("POST")
	r.HandleFunc("/charges/{id}", getChargeHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	log.Println("PayPal mock server running on :8081")
	http.ListenAndServe(":8081", r)
}
