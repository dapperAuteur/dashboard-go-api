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
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjEiLCJ0eXAiOiJKV1QifQ.eyJyb2xlcyI6WyJBRE1JTiIsIlVTRVIiXSwiZXhwIjoxNTk4MDY1MTAzLCJpYXQiOjE1OTgwNjE1MDMsInN1YiI6IjVmMzMyZjFjNzBiMmVjMzIwYmMyMTlkMCJ9.u--ffPU_-GAbWVCxF1tj_PShwg9cvlh7VT4jGx9SdG1vv-x_a0cQBROEujGBwhY1hIE_dnWgEDhoY-SuNWiXjSCE99Ts_PA3wui_ivWPUnc9RaOno75K_HN8e8SmrRKOjQzpPTNuFl2CN2jNh5AofptSw6SEIplmtnybY8y7NZdky4Ys_w7Zw8vAo29pOnCq9rERYxx71SiAdGJqFycK8bac0k8xd6FR_cCvEl1n-CV7lVRedKuIyKD4XDkcgoHLAU_-6IV-XlJ_rqUCI595Vj_BgfP8b9UnWgMTuSIKb1yn0bTBNYkDDDA9FSy1Zy-UP9KecNTILlKJ92CDqkzPOA")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
