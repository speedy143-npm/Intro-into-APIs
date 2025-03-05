package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"

	"log"
	"net/http"
)

// Struct types to hold output
type Transreq struct {
	From        string `json:"from"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Reference   string `json:"external_reference"`
}

type Transrep struct {
	Reference string `json:"reference"`
	//Status    string `json:"status"`
	Ussd_code string `json:"ussd_code"`
	Operator  string `json:"operator"`
}

type State struct {
	Reference string `json:"reference"`
	Ext_ref   string `json:"external_reference"`
	Status    string `json:"status"`
	// "amount"
	// "currency"
	// "operator"
	// "code"
	// "operator_reference"
	// "description"
	// "external_user"
	// "reason"
	// "phone_number"
}

func main() {
	var number string
	var amount string
	var description, ref string

	apikey, err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Requesting inputs from user
	for {
		fmt.Println("Enter your mobile money number without country code")
		fmt.Scanln(&number)
		if isValidnumber(number) {
			break
		}
		fmt.Println("Invalid phone number. Please enter a valid phone number starting with 675, 673, 651, 678 or 677 followed by exactly 6 other numbers.")
	}

	number = "237" + number

	fmt.Println("Enter amount")
	fmt.Scanln(&amount)

	fmt.Println("Enter description")
	fmt.Scanln(&description)

	fmt.Println("Enter Reference")
	fmt.Scanln(&ref)

	trans := post(apikey, number, amount, description, ref)
	fmt.Printf("Transaction Reference: %s\nTransaction Code: %s\n", trans.Reference, trans.Ussd_code)

	//waiting time before checking transaction status
	time.Sleep(30 * time.Second)

	state := get(apikey, trans.Reference)
	fmt.Printf("Status: %s\n", state.Status)

}

// Function to validate phone number
func isValidnumber(number string) bool {
	re := regexp.MustCompile(`^(675|673|651|677|678)\d{6}$`)
	return re.MatchString(number)
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

// Initiating mobile money withdrawal
func post(apik string, number string, amount string, description string, ref string) Transrep {

	client := &http.Client{}

	transreq := Transreq{
		From:        number,
		Amount:      amount,
		Description: description,
		Reference:   ref,
	}

	reqBody, _ := json.Marshal(transreq)

	req, err := http.NewRequest("POST", "https://demo.campay.net/api/collect/", bytes.NewBuffer(reqBody))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apik))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var transaction Transrep
	json.NewDecoder(response.Body).Decode(&transaction)
	return transaction

}

// Initiating request to get the status of the transaction
func get(apik, reference string) State {
	client := &http.Client{}

	url := fmt.Sprintf("https://demo.campay.net/api/transaction/%s", reference)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apik))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var checkState State
	json.NewDecoder(response.Body).Decode((&checkState))
	return checkState

}
