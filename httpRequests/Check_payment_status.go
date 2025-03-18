package httpRequests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// struct type to hold output
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

// Initiating request to get the status of the transaction
func CheckPaymentStatus(apik, reference string) Status {
	client := &http.Client{}

	url := fmt.Sprintf("https://demo.campay.net/api/transaction/%s", reference)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Check GET request credentials")
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apik))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)

	if err != nil {
		fmt.Println("Invalid Request, check post request credentials")
		log.Fatal(err)
	}

	defer response.Body.Close()

	var checkState Status
	json.NewDecoder(response.Body).Decode((&checkState))
	return checkState

}
