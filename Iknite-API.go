package main

import (
	"fmt"
	"os"

	"time"

	"iknite-api/httpRequests"

	"github.com/joho/godotenv"
)

func main() {
	var number string
	var amount string
	var description, ref string

	apikey, err := run()
	if err != nil {
		fmt.Println("API KEY could not be found\n", err)
		os.Exit(1)
	}

	// Requesting inputs from user

	trans := httpRequests.RequestPayment(apikey, number, amount, description, ref)
	fmt.Printf("Transaction Reference: %s\nTransaction Code: %s\n", trans.Reference, trans.Ussd_code)

	//waiting time before checking transaction status
	time.Sleep(30 * time.Second)

	state := httpRequests.CheckPaymentStatus(apikey, trans.Reference)
	fmt.Printf("Status: %s\n", state.Status)

}

// loading the API KEY from .env file
func run() (string, error) {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			return "", fmt.Errorf("failed to load env file: %w", err)
		}
	}
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("API_KEY is not set")
	}

	return apiKey, nil

}
