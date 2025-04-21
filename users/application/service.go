package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/MikalaiLappo/user_billing_micros/users/domain"
	"github.com/google/uuid"
)

type UserService struct {
	billingMicroURL string
}

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

func NewUserService(billingUrl string) *UserService {
	return &UserService{
		billingMicroURL: billingUrl,
	}
}

func (s *UserService) RunSpeculationThread(ctx context.Context, id uuid.UUID, secFreq uint32) error {
	ticker := time.NewTicker(time.Second * time.Duration(secFreq))
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// 50/50
				if boolRand() {
					// recieve case
					fmt.Printf("[RunSpeculationThread] making a request\n")
					data := []byte(fmt.Sprintf(`{"owner_id": "%s", "amount": %d}`, id.String(), 100))
					r := bytes.NewReader(data)

					resp, err := http.Post(s.billingMicroURL+"/recieve", "application/json", r)
					if err != nil {
						fmt.Printf("[RunSpeculationThread] failed to recieve balance user_id: %s | %s\n", id.String(), err.Error())
						continue
					}

					if resp.StatusCode != http.StatusOK {
						fmt.Printf("[RunSpeculationThread] billing service coudldn't perform a recieve operation\n")
					}

					fmt.Printf("[RunSpeculationThread] recieve: OK\n")
					continue
				}
				// pay case
				fmt.Printf("[RunSpeculationThread] making a request\n")
				data := []byte(fmt.Sprintf(`{"owner_id": "%s", "amount": %d}`, id.String(), 100))
				r := bytes.NewReader(data)

				resp, err := http.Post(s.billingMicroURL+"/pay", "application/json", r)
				if err != nil {
					fmt.Printf("[RunSpeculationThread] failed to pay balance user_id: %s | %s\n", id.String(), err.Error())
					continue
				}

				if resp.StatusCode != http.StatusOK {
					fmt.Printf("[RunSpeculationThread] billing service coudldn't make a pay\n")
				}

				fmt.Printf("[RunSpeculationThread] pay: OK\n")

			case <-quit:
				ticker.Stop()
				fmt.Println("[RunSpeculationThread] quit\n")
			}
		}
	}()

	return nil
}

func (s *UserService) CheckBalanceThread(ctx context.Context, id uuid.UUID, secFreq uint32) error {
	ticker := time.NewTicker(time.Second * time.Duration(secFreq))
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Printf("[CheckBalanceThread] making a request\n")
				data := []byte(fmt.Sprintf(`{"owner_id": "%s"}`, id.String()))
				r := bytes.NewReader(data)

				_, err := http.Post(s.billingMicroURL+"/balance", "application/json", r)
				if err != nil {
					fmt.Printf("[CheckBalanceThread] failed to retrieve balance user_id: %s | %s\n", id.String(), err.Error())
					continue
				}
				fmt.Printf("[CheckBalanceThread] balance response: OK\n")

			case <-quit:
				ticker.Stop()
				fmt.Println("[CheckBalanceThread] quit")
			}
		}
	}()

	return nil
}

func PrintStruct(v interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error converting struct to JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func boolRand() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
