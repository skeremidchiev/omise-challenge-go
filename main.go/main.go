package main

import (
	"errors"
	"bufio"
	"os"
	"fmt"
	"log"
	"omise_challenges/cipher"
)

func main () {
	// TODO: move this to separete functio
	// TODO: add file path as command-line option
	file, err := os.Open("../data/fng.1000.csv.rot128")
	if err != nil {
		log.Fatal(err)
	}

	defer func () {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	stats, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)
	rot12Reader, err := cipher.NewRot128Reader(reader)
	buffer := make([]byte, stats.Size())

	n, err := rot12Reader.Read(buffer)
	// EOF error won't be reached
	if err != nil {
		log.Fatal(err)
	}

	if int64(n) < stats.Size() {
		log.Fatal(errors.New("File not read fully"))
	}

	// get buffer
	log.Println(string(buffer[0:n]))
	// clear buffer
	buffer = nil
	log.Println("after nil:", buffer)
}