package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MikalaiLappo/user_billing_micros/billing/application"
	"github.com/MikalaiLappo/user_billing_micros/billing/domain"
	"github.com/MikalaiLappo/user_billing_micros/billing/infrastructure"
	"github.com/MikalaiLappo/user_billing_micros/billing/interfaces/rest"
	"github.com/google/uuid"
)

func main() {

	billingRepo := infrastructure.NewInMemoryBillingRepository()

	// mock
	mockAccount := domain.NewBillingAccount()
	hardcodedId, _ := uuid.Parse("b22084ec-6da9-4767-8eb9-5a4bfa9f3a37")
	mockAccount.OwnerID = hardcodedId
	fmt.Printf("[MOCK] billing account for testing: %+v\n", mockAccount)
	billingRepo.Save(context.Background(), mockAccount)
	//

	billingService := application.NewBillingService(billingRepo)
	billingHandler := rest.NewBillingHandler(billingService)
	http.HandleFunc("/balance", billingHandler.GetAccountBalance)
	http.HandleFunc("/pay", billingHandler.PayWithAccount)
	http.HandleFunc("/recieve", billingHandler.RecieveMoney)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err.Error())
	}

	select {}

}
