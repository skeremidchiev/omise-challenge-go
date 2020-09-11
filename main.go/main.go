package main

import (
	"fmt"
	"omise_challenges/parser"
)

func main () {
	users := make(chan string)

	go func () {
		parser.Parse("../data/fng.1000.csv.rot128", users)
	}()

	for user := range users {
		fmt.Println(user)
	}
}