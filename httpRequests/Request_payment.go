package httpRequests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type Requests struct {
	BaseUrl string
	Apikey  string
}

func NewApiClient(baseurl string) *Requests {
	return &Requests{
		BaseUrl: baseurl,
		// Apikey: apikey,
	}
}

// Struct types to hold input
type Transrequest struct {
	From        string `json:"from"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Reference   string `json:"external_reference"`
}

// Struct types to hold output
type Transresponse struct {
	Reference string `json:"reference"`
	Ussd_code string `json:"ussd_code"`
	Operator  string `json:"operator"`
}

// Function to validate phone number
func IsValidnumber(number string) bool {
	re := regexp.MustCompile(`^(670|671|672|673|674|675|676|677|678|679|680|681|682|683|684|685|686|687|688|689|650|651|652|653|654|691|692|693|694|694|695|696|696|697|698|699|655|656|657|658|659|620|621|622|623)\d{6}$`)
	return re.MatchString(number)
}

func IsValidamount(amount string) bool {
	re := regexp.MustCompile(`^([1-9][0-9]{0,4}|500000)$`)
	return re.MatchString(amount)
}

// Initiating mobile money withdrawal
func RequestPayment(apik string, number string, amount string, description string, ref string) Transresponse {

	//Requesting inputs from user
	for {
		fmt.Println("Enter your mobile money number without country code")
		fmt.Scanln(&number)
		if IsValidnumber(number) {
			break
		}
		fmt.Println("Invalid phone number. Please enter a valid phone number starting with 675, 673, 651, 653, 680, 678 or 677 followed by exactly 6 other numbers.")
	}

	number = "237" + number
	for {
		fmt.Println("Enter Amount ")
		fmt.Scanln(&amount)
		if IsValidamount(amount) {
			break
		}
		fmt.Println("Invalid amount. Please enter an amount within 1 and 500,000")
	}

	fmt.Println("Enter description")
	fmt.Scanln(&description)

	fmt.Println("Enter Reference")
	fmt.Scanln(&ref)

	//Creating a http client
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
