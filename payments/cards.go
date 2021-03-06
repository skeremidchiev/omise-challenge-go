package payments

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"omise_challenges/config"
)

type CardToken struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

type CCard struct {
	Name     string `json:"name"`
	Number   string `json:"number"`
	ExpMonth string `json:"expiration_month"`
	ExpYear  string `json:"expiration_year"`
	SecCode  string `json:"security_code"`
}

// This piece of code is from omise-go/operations/token.go
func (c *CCard) MarshalJSON() ([]byte, error) {
	type CardToken CCard
	params := struct {
		Card *CardToken `json:"card"`
	}{
		Card: (*CardToken)(c),
	}
	return json.Marshal(params)
}

func (c *CCard) CreateToken(token *CardToken) error {
	config := config.GetConfig()

	data, err := c.MarshalJSON()
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		"POST",
		config.TokensUrl,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(config.PubKey, "")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// In all cases of StatusCode == 400 it's expired card
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"Request failed with error code: %d",
			response.StatusCode,
		)
	}

	err = json.NewDecoder(response.Body).Decode(token)
	if err != nil {
		return err
	}

	return nil
}