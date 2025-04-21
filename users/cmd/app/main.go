package main

import (
	"context"
	"fmt"

	"github.com/MikalaiLappo/user_billing_micros/users/application"
	"github.com/MikalaiLappo/user_billing_micros/users/domain"
	"github.com/MikalaiLappo/user_billing_micros/users/infrastructure"
	"github.com/google/uuid"
)

func main() {
	userRepo := infrastructure.NewInMemoryBillingRepository()
	// mock
	mockAccount := domain.NewUser()
	hardcodedId, _ := uuid.Parse("b22084ec-6da9-4767-8eb9-5a4bfa9f3a37")
	mockAccount.ID = hardcodedId
	fmt.Printf("[MOCK] billing account for testing: %+v\n", mockAccount)
	userRepo.Save(context.Background(), mockAccount)
	//

	userService := application.NewUserService(userRepo, "http://localhost:8000")

	ctx := context.Background()
	userService.RunSpeculationThread(ctx, hardcodedId, 1)
	go func() { userService.CheckBalanceThread(context.Background(), hardcodedId, 1) }()

	select {}
}
