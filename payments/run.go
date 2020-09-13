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

func worker(jobsChan <-chan string, resultChan chan<- *PaymentResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range jobsChan {
		// fmt.Printf("USER:\n%s\n", data)

		dataArr, valid := parser.SplitAndValidate(data)
		if !valid {
			// dropping invalid user data
			continue
		}

		result := Pay(dataArr)
		// fmt.Println(result)
		resultChan <- result
	}
}

func Run(donationsFP, configFP string) {
	fmt.Println("performing donations...")

	// initial set up
	var wg sync.WaitGroup
	jobsChan := make(chan string) // listen for a new job
	resultChan := make(chan *PaymentResult, WORKERS_LIMIT * 2) // used for summary
	sum := GetSummaryObject()
	config.SetConfigFilePath(configFP)

	// I'm hitting the rate limits with more than 3 go routines
	// 429 - Too many requests and 500 - ???
	// This needs more dynamic and adaptable approach,
	// mine is too static and not 100% effective
	for w := 1; w <= WORKERS_LIMIT; w++ {
		wg.Add(1)
		go worker(jobsChan, resultChan, &wg)
	}

	// summary
	go func() {
		for r := range resultChan {
			sum.Add(r)
		}
	}()

	// feed jobs channel
	go func () {
		err := parser.Parse(donationsFP, jobsChan)
		if err != nil {
			log.Println("Error During Parsing: ", err)
		}
		// Close the output channel when done
		defer close(jobsChan)
	}()

	wg.Wait()
	fmt.Println("done.")
	fmt.Println(sum.GetSummary())
}