package payments

import (
	"log"
	"strconv"
)

type PaymentResult struct {
	Name   string
	Amount int
	Status bool
}

func Pay(userData []string) *PaymentResult {
	donationAmount, _ := strconv.Atoi(userData[1])
	result := &PaymentResult {
		userData[0],
		donationAmount,
		false,
	}

	cardToken, card := &CardToken{}, &CCard{
		Name:     userData[0],
		Number:   userData[2],
		ExpMonth: userData[4],
		ExpYear:  userData[5],
		SecCode:  userData[3],
	}

	// Create card token
	err := card.CreateToken(cardToken)
	if err != nil {
		log.Println("Error During Token Creation: ", err)
		return result
	}

	chargeResult, charge := &ChargeResult{}, &ChargeData{
		Token: cardToken.ID,
		Currency: DEFAULT_CURRENCY,
		Amount: userData[1],
	}

	// charge
	err = charge.Charge(chargeResult)
	if err != nil {
		log.Println("Error During Charging: ", err)
		return result
	}

	if chargeResult.Status != "successful" {
		return result
	}

	result.Status = true
	return result
}
