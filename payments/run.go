package payments

import (
	"fmt"
	"log"
	"sync"
	"omise_challenges/parser"
)

const (
	WORKERS_LIMIT = 3
)
// "sync"
// var mutex = &sync.Mutex{}
// mutex.Lock()
// mutex.Unlock()

func worker(jobsChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range jobsChan {
		fmt.Printf("USER:\n%s\n", data)

		dataArr, valid := parser.SplitAndValidate(data)
		if !valid {
			// not handling invalid user data
			continue
		}

		result := Pay(dataArr)
		fmt.Println(result)
	}
}

func Run() {
	var wg sync.WaitGroup

	jobsChan := make(chan string) // listen for a new job

	// I'm hitting the rate limits with more than 3 go routines
	// 429 - Too many requests and 500 - ???
	// This needs more dynamic and adaptable approach,
	// mine is too static and not 100% effective
	for w := 1; w <= WORKERS_LIMIT; w++ {
		wg.Add(1)
		go worker(jobsChan, &wg)
	}

	// feed jobs channel
	go func () {
		err := parser.Parse("../data/fng.1000.csv.rot128", jobsChan)
		if err != nil {
			log.Println("Error During Parsing: ", err)
		}
		// Close the output channel when done
		defer close(jobsChan)
	}()

	wg.Wait()
}