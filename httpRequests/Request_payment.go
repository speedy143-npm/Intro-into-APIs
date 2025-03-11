package httpRequests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

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

// Function to validate phone number
func IsValidnumber(number string) bool {
	re := regexp.MustCompile(`^(675|673|651|677|678|653|680)\d{6}$`)
	return re.MatchString(number)
}

// Initiating mobile money withdrawal
func RequestPayment(apik string, number string, amount string, description string, ref string) Transresponse {

	client := &http.Client{}

	transreq := Transrequest{
		From:        number,
		Amount:      amount,
		Description: description,
		Reference:   ref,
	}

	reqBody, _ := json.Marshal(transreq)

	req, err := http.NewRequest("POST", "https://demo.campay.net/api/collect/", bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Println("Check post request credentials")
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apik))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Invalid Request, check POST request credentials")
		log.Fatal(err)
	}

	defer response.Body.Close()

	var transaction Transresponse
	json.NewDecoder(response.Body).Decode(&transaction)
	return transaction

}
