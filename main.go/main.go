package main

import (
	"fmt"
	"log"
	"omise_challenges/parser"
	"omise_challenges/payments"
)

func main () {
	users := make(chan string)

	go func () {
		err := parser.Parse("../data/fng.1000.csv.rot128", users)
		if err != nil {
			log.Println("Error During Parsing: ", err)
		}
	}()

	for user := range users {
		fmt.Printf("USER:\n%s\n", user)

		dataArr, valid := parser.SplitAndValidate(user)
		if !valid {
			continue
		}

		cardToken, card := &payments.CardToken{}, &payments.CCard{
			Name:     dataArr[0],
			Number:   dataArr[2],
			ExpMonth: dataArr[4],
			ExpYear:  dataArr[5],
			SecCode:  dataArr[3],
		}

		err := card.CreateToken(cardToken)
		if err != nil {
			log.Println("Error During Token Creation: ", err)
			continue
		}

		fmt.Println(cardToken)
	}
}