package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MikalaiLappo/user_billing_micros/billing/application"
	"github.com/MikalaiLappo/user_billing_micros/billing/interfaces/rest/dto"
	"github.com/google/uuid"
)

type BillingHandler struct {
	service *application.BillingService
}

func NewBillingHandler(service *application.BillingService) *BillingHandler {
	return &BillingHandler{
		service: service,
	}
}

func (h *BillingHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GetAccountBalance] | requested | %s\n", time.Now().Format(time.RFC3339))
	w.Header().Set("Context-Type", "application/json")

	var reqPayload dto.GetBalanceRequest
	err := json.NewDecoder(r.Body).Decode(&reqPayload)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[GetAccountBalance] | error (invalid user request) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	uuid, err := uuid.Parse(reqPayload.OwnerID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[GetAccountBalance] | error (invalid user_id format) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	balance, err := h.service.GetBillingBalance(r.Context(), uuid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[GetAccountBalance] | error (repo) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetBalanceResponse{
		OwnerID: uuid.String(),
		Balance: balance,
	})

	log.Println("[GetAccountBalance] OK")
}

func (h *BillingHandler) PayWithAccount(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PayWithAccount] | requested | %s\n", time.Now().Format(time.RFC3339))
	w.Header().Set("Context-Type", "application/json")

	var reqPayload dto.PayRequest
	err := json.NewDecoder(r.Body).Decode(&reqPayload)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[PayWithAccount] | error (invalid user request) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	uuid, err := uuid.Parse(reqPayload.OwnerID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[PayWithAccount] | error (invalid user_id format) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	err = h.service.Pay(r.Context(), uuid, reqPayload.Amount)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[PayWithAccount] | error (repo pay) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	balance, err := h.service.GetBillingBalance(r.Context(), uuid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[PayWithAccount] | error (repo get) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.PayResponse{
		OwnerID:    uuid.String(),
		NewBalance: balance,
	})

	log.Println("[PayWithAccount] OK")
}

func (h *BillingHandler) RecieveMoney(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RecieveMoney] | requested | %s\n", time.Now().Format(time.RFC3339))
	w.Header().Set("Context-Type", "application/json")

	var reqPayload dto.RecieveRequest
	err := json.NewDecoder(r.Body).Decode(&reqPayload)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[RecieveMoney] | error (invalid user request) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	uuid, err := uuid.Parse(reqPayload.OwnerID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[RecieveMoney] | error (invalid user_id format) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	err = h.service.RecieveMoney(r.Context(), uuid, reqPayload.Amount)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[RecieveMoney] | error (repo pay) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	balance, err := h.service.GetBillingBalance(r.Context(), uuid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errFmt := fmt.Sprintf("[PayWithAccount] | error (repo get) | %s\n", err.Error())
		log.Println(errFmt)
		w.Write([]byte(errFmt))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.RecieveResponse{
		OwnerID:    uuid.String(),
		NewBalance: balance,
	})

	log.Println("[RecieveMoney] OK")
}
