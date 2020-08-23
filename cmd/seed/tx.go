package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	url := "http://localhost:8080/v1/transactions"
	method := "POST"

	payload := strings.NewReader("{\n      \"participant_id\": [\"5ab055eae67be20014ca5284\"],\n      \"vendor_id\": \"5f3e18f8d95d06627dc8e968\",\n      \"tranx_event\": \"groceries\",\n      \"occurrence_string\": \"8/20\",\n      \"budget_id\": \"5f3e189bd95d06627dc8e932\",\n      \"currency_id\": \"5f381f30f815d062fb9da8f1\",\n      \"fin_acc_id\": [\"5f3e16a8d95d06627dc8e92f\"],\n      \"tranx_debit\": 352,\n      \"tranx_credit\": 176\n    }")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjEiLCJ0eXAiOiJKV1QifQ.eyJyb2xlcyI6WyJBRE1JTiIsIlVTRVIiXSwiZXhwIjoxNTk4MDc3NTE4LCJpYXQiOjE1OTgwNzM5MTgsInN1YiI6IjVmMzMyZjFjNzBiMmVjMzIwYmMyMTlkMCJ9.M9vTuWDTeHHMAvxvvU6M1lgralNGRhtnwqUITbd4UlaWd21-rZCyUKXDFpS_LVNZBpe5LbimtE6wGMpreU-W1zpu7u-7125cqngg7Is8TW5uyZg4_Csg0oEy5UrtTEC0DY8-rfRol4pN2nmqvp_CQ_uaiK_FVA-LpM4xPT-d6_VpLwSzUJPINMRJNC_VQu87S3TXGjSOisqtd5uFH1DS6aopwen_T8hHwd7go1tmixPAHnRUFKp51zuJe0irwCiJIjQ8SUarwdu-DDjKh5QJ11w-PU_DG1UjcZGevwB1voOwhbAf6Y7dpIAr4yjs-5lXflGKBRiKkhq8rrZH20bFZg")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
