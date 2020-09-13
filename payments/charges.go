package payments

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"omise_challenges/config"
)

const (
	DEFAULT_CURRENCY = "thb"
)

type ChargeResult struct {
	Charge string `json:"id"`
	Status string `json:"status"`
}

type ChargeData struct {
	Token       string `json:"card"`
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
}

func (pd *ChargeData) Charge(result *ChargeResult) error {
	config := config.GetConfig()

	data, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		"POST",
		config.ChargeUrl,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(config.SecKey, "")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"Request failed with error code: %d",
			response.StatusCode,
		)
	}

	data, _ = ioutil.ReadAll(response.Body)
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return err
	}

	return nil
}