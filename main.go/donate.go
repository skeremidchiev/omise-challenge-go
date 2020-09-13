package main

import (
	"flag"
	"omise_challenges/payments"
)

func main () {
	donateFilePathFlag := flag.String("dfpath", "", "Donation file to parse.")
	confingFilePathFlag := flag.String("cfpath", "", "Configuration file.")
	flag.Parse()

	payments.Run(*donateFilePathFlag, *confingFilePathFlag)
}