package payments

import (
	"fmt"
	"log"
	"sync"
	"omise_challenges/parser"
	"omise_challenges/config"
)

const (
	WORKERS_LIMIT = 3
)

func worker(jobsChan <-chan string, sum *Summary, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range jobsChan {
		dataArr, valid := parser.SplitAndValidate(data)
		if !valid {
			// dropping invalid user data
			continue
		}

		result := Pay(dataArr)

		// adding data to summary
		wg.Add(1)
		go func(r *PaymentResult) {
			defer wg.Done()
			sum.Add(r)
		} (result)
	}
}

func Run(donationsFP, configFP string) {
	fmt.Println("performing donations...")

	// initial set up
	config.SetConfigFilePath(configFP)
	var wg sync.WaitGroup
	jobsChan := make(chan string) // listen for a new job
	sum := GetSummaryObject()

	// I'm hitting the rate limits with more than 3 go routines
	// 429 - Too many requests and 500 - ???
	// This needs more dynamic and adaptable approach,
	// mine is too static and not 100% effective
	for w := 1; w <= WORKERS_LIMIT; w++ {
		wg.Add(1)
		go worker(jobsChan, sum, &wg)
	}

	// feed jobs channel
	go func () {
		err := parser.Parse(donationsFP, jobsChan)
		if err != nil {
			log.Println("Error During Parsing: ", err)
		}
		// Close the output channel when done parsing
		defer close(jobsChan)
	}()

	wg.Wait()
	fmt.Println("done.")
	fmt.Println(sum.GetSummary())
}