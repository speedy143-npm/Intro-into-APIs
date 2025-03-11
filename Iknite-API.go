package main

import (
	"fmt"
	"os"

	"time"

	"iknite-api/httpRequests"

	"github.com/joho/godotenv"
)

// Storing RequestPayment and CheckPaymentStatus requests on struct
type Requests struct {
	RequestPayment     string
	CheckPaymentStatus string
	BaseUrl            string `json:"https://demo.campay.net"`
	ApiKey             string
}

// Struct types to hold output
type Transrequest struct {
	From        string `json:"from"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Reference   string `json:"external_reference"`
}

type Transresponse struct {
	Reference string `json:"reference"`
	//Status    string `json:"status"`
	Ussd_code string `json:"ussd_code"`
	Operator  string `json:"operator"`
}

type Status struct {
	Reference          string `json:"reference"`
	Ext_ref            string `json:"external_reference"`
	Status             string `json:"status"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	Operator           string `json:"operator"`
	Code               string `json:"code"`
	Operator_Reference string `json:"operator_reference"`
	Description        string `json:"description"`
	Exterbal_User      string `json:"external_user"`
	Reason             string `json:"reason"`
	Phone_Number       string `json:"phone_number"`
}

func main() {
	x := httpRequests.Add(1, 2)
	fmt.Println(x)
	var number string
	var amount string
	var description, ref string

	apikey, err := run()
	if err != nil {
		fmt.Println("API KEY could not be found\n", err)
		os.Exit(1)
	}

	// Requesting inputs from user
	for {
		fmt.Println("Enter your mobile money number without country code")
		fmt.Scanln(&number)
		if httpRequests.IsValidnumber(number) {
			break
		}
		fmt.Println("Invalid phone number. Please enter a valid phone number starting with 675, 673, 651, 653, 680, 678 or 677 followed by exactly 6 other numbers.")
	}

	number = "237" + number

	fmt.Println("Enter amount")
	fmt.Scanln(&amount)

	fmt.Println("Enter description")
	fmt.Scanln(&description)

	fmt.Println("Enter Reference")
	fmt.Scanln(&ref)

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
