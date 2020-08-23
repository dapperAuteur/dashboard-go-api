package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	// Open json file
	// jsonFile, err := os.Open("tranx_02.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer jsonFile.Close()
	// fmt.Println("file opened successfully")

	// initialize NewTransaction
	// var newTranxSlice []NewTransaction
	// var newTranx NewTransaction

	// // NewTransaction type is what's required from the client to create a new transaction.
	// type NewTransaction struct {
	// 	BudgetID           primitive.ObjectID `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	// 	CurrencyID         primitive.ObjectID `bson:"currency_id,omitempty" json:"currency_id,omitempty"`
	// 	FinancialAccountID *[]string          `bson:"fin_acc_id,omitempty" json:"fin_acc_id,omitempty"`
	// 	// Occurrence         *time.Time            `bson:"occurrence,omitempty" json:"occurrence,omitempty" validate:"datetime"`
	// 	OccurrenceString  string             `bson:"occurrence_string,omitempty" json:"occurrence_string,omitempty"`
	// 	TransactionEvent  string             `bson:"tranx_event,omitempty" json:"tranx_event,omitempty"`
	// 	TransactionCredit float64            `bson:"tranx_credit,omitempty" json:"tranx_credit,omitempty"`
	// 	TransactionDebit  float64            `bson:"tranx_debit,omitempty" json:"tranx_debit,omitempty"`
	// 	VendorID          primitive.ObjectID `bson:"vendor_id,omitempty" json:"vendor_id,omitempty"`
	// 	ParticipantID     *[]string          `bson:"participant_id,omitempty" json:"participant_id,omitempty"`
	// }

	Data := []byte(`
	[
		{
			"participant_id": ["5ab055eae67be20014ca5284"],
			"vendor_id": "5f3e18f8d95d06627dc8e94e",
			"tranx_event": "movies-Bad Boy For Life",
			"occurrence_string": "3/1",
			"budget_id": "5f3e189bd95d06627dc8e931",
			"currency_id": "5f381f30f815d062fb9da8f1",
			"fin_acc_id": ["5f3e16a8d95d06627dc8e928"],
			"tranx_debit": 1325,
			"tranx_credit": 0
		  },
		  {
			"participant_id": ["5ab055eae67be20014ca5284"],
			"vendor_id": "5f3e18f8d95d06627dc8e990",
			"tranx_event": "groceries: alkalizer & detoxifier supplement",
			"occurrence_string": "3/1",
			"budget_id": "5f3e189bd95d06627dc8e932",
			"currency_id": "5f381f30f815d062fb9da8f1",
			"fin_acc_id": ["5f3e16a8d95d06627dc8e928"],
			"tranx_debit": 2928,
			"tranx_credit": 0
		  },
		  {
			"participant_id": ["5ab055eae67be20014ca5284"],
			"vendor_id": "5f3e18f8d95d06627dc8e952",
			"tranx_event": "tire gauge 120 psi",
			"occurrence_string": "3/1",
			"budget_id": "5f3e189bd95d06627dc8e937",
			"currency_id": "5f381f30f815d062fb9da8f1",
			"fin_acc_id": ["5f3e16a8d95d06627dc8e928"],
			"tranx_debit": 720,
			"tranx_credit": 0
		  }
	]
	`)
	// unmarshal byteArray which contains
	// jsonFile's content into '[]NewTransaction' which we defined above
	var tx []NewTransaction
	err := json.Unmarshal(Data, &tx)
	// err := json.Unmarshal([]byte(Data), &tx)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("****************\ntx : ", tx)
	fmt.Println("****************\nreflect.TypeOf(tx).String() : \n", reflect.TypeOf(tx).String())

	// fmt.Println("Tranx : ", tx)

	// for _, t := range tx {
	// 	fmt.Println("Tranx : ", t.TransactionDebit)
	// }

	// fmt.Println("Tx : ", tx)

	// for i := 0; i < len(tx.Tx); i++ {
	// 	fmt.Println("Tx : ", tx.Tx[i])
	// 	fmt.Println("Tx : ", tx.Tx[i].TransactionDebit)

	// }

	// fmt.Println("Tranx : ", &newTranxSlice)
	url := "http://localhost:8080/v1/transactions"
	method := "POST"

	// make POST request
	for _, t := range tx {
		fmt.Println("Tranx : ", t.TransactionDebit)
		// payload := strings.NewReader("{\n      \"participant_id\": [\"5ab055eae67be20014ca5284\"],\n      \"vendor_id\": \"5f3e18f8d95d06627dc8e968\",\n      \"tranx_event\": \"groceries\",\n      \"occurrence_string\": \"8/20\",\n      \"budget_id\": \"5f3e189bd95d06627dc8e932\",\n      \"currency_id\": \"5f381f30f815d062fb9da8f1\",\n      \"fin_acc_id\": [\"5f3e16a8d95d06627dc8e92f\"],\n      \"tranx_debit\": 352,\n      \"tranx_credit\": 176\n    }")

		// nTx, err := json.Unmarshal()
		payload := strings.NewReader()

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjEiLCJ0eXAiOiJKV1QifQ.eyJyb2xlcyI6WyJBRE1JTiIsIlVTRVIiXSwiZXhwIjoxNTk4MDY1MTAzLCJpYXQiOjE1OTgwNjE1MDMsInN1YiI6IjVmMzMyZjFjNzBiMmVjMzIwYmMyMTlkMCJ9.u--ffPU_-GAbWVCxF1tj_PShwg9cvlh7VT4jGx9SdG1vv-x_a0cQBROEujGBwhY1hIE_dnWgEDhoY-SuNWiXjSCE99Ts_PA3wui_ivWPUnc9RaOno75K_HN8e8SmrRKOjQzpPTNuFl2CN2jNh5AofptSw6SEIplmtnybY8y7NZdky4Ys_w7Zw8vAo29pOnCq9rERYxx71SiAdGJqFycK8bac0k8xd6FR_cCvEl1n-CV7lVRedKuIyKD4XDkcgoHLAU_-6IV-XlJ_rqUCI595Vj_BgfP8b9UnWgMTuSIKb1yn0bTBNYkDDDA9FSy1Zy-UP9KecNTILlKJ92CDqkzPOA")

		res, err := client.Do(req)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		fmt.Println(string(body))
	}

	// payload := strings.NewReader("{\n      \"participant_id\": [\"5ab055eae67be20014ca5284\"],\n      \"vendor_id\": \"5f3e18f8d95d06627dc8e968\",\n      \"tranx_event\": \"groceries\",\n      \"occurrence_string\": \"8/20\",\n      \"budget_id\": \"5f3e189bd95d06627dc8e932\",\n      \"currency_id\": \"5f381f30f815d062fb9da8f1\",\n      \"fin_acc_id\": [\"5f3e16a8d95d06627dc8e92f\"],\n      \"tranx_debit\": 352,\n      \"tranx_credit\": 176\n    }")

	// client := &http.Client{}
	// req, err := http.NewRequest(method, url, payload)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjEiLCJ0eXAiOiJKV1QifQ.eyJyb2xlcyI6WyJBRE1JTiIsIlVTRVIiXSwiZXhwIjoxNTk4MDY1MTAzLCJpYXQiOjE1OTgwNjE1MDMsInN1YiI6IjVmMzMyZjFjNzBiMmVjMzIwYmMyMTlkMCJ9.u--ffPU_-GAbWVCxF1tj_PShwg9cvlh7VT4jGx9SdG1vv-x_a0cQBROEujGBwhY1hIE_dnWgEDhoY-SuNWiXjSCE99Ts_PA3wui_ivWPUnc9RaOno75K_HN8e8SmrRKOjQzpPTNuFl2CN2jNh5AofptSw6SEIplmtnybY8y7NZdky4Ys_w7Zw8vAo29pOnCq9rERYxx71SiAdGJqFycK8bac0k8xd6FR_cCvEl1n-CV7lVRedKuIyKD4XDkcgoHLAU_-6IV-XlJ_rqUCI595Vj_BgfP8b9UnWgMTuSIKb1yn0bTBNYkDDDA9FSy1Zy-UP9KecNTILlKJ92CDqkzPOA")

	// res, err := client.Do(req)
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)

	// fmt.Println(string(body))

}

// NewTransaction type is what's required from the client to create a new transaction.
type NewTransaction struct {
	BudgetID           primitive.ObjectID `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	CurrencyID         primitive.ObjectID `bson:"currency_id,omitempty" json:"currency_id,omitempty"`
	FinancialAccountID *[]string          `bson:"fin_acc_id,omitempty" json:"fin_acc_id,omitempty"`
	// Occurrence         *time.Time            `bson:"occurrence,omitempty" json:"occurrence,omitempty" validate:"datetime"`
	OccurrenceString  string             `bson:"occurrence_string,omitempty" json:"occurrence_string,omitempty"`
	TransactionEvent  string             `bson:"tranx_event,omitempty" json:"tranx_event,omitempty"`
	TransactionCredit float64            `bson:"tranx_credit,omitempty" json:"tranx_credit,omitempty"`
	TransactionDebit  float64            `bson:"tranx_debit,omitempty" json:"tranx_debit,omitempty"`
	VendorID          primitive.ObjectID `bson:"vendor_id,omitempty" json:"vendor_id,omitempty"`
	ParticipantID     *[]string          `bson:"participant_id,omitempty" json:"participant_id,omitempty"`
}

func createTranxHandler(method, pattern string, h http.ResponseWriter, r *http.Request) {

}
