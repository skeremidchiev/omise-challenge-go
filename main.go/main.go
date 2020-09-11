package main

import (
	// "errors"
	"strings"
	"bufio"
	"os"
	// "fmt"
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

	reader := bufio.NewReader(file)
	rot12Reader, err := cipher.NewRot128Reader(reader)

	const INITIAL_SIZE = 50
	buffer := make([]byte, INITIAL_SIZE)

	start := 0
	for {
		n, err := rot12Reader.Read(buffer[start:])
		if err != nil {
			log.Fatal(err)
		}
		log.Println("BYTES READ: ", n)
		// get buffer
		log.Println("ORIGINAL: \n",string(buffer))

		lastIndexOfNewLine := strings.LastIndex(string(buffer), "\n")
		log.Println("Last index of :", lastIndexOfNewLine)

		// In this case extend the buffer capacity
		// Try to keep buffer capacity as low as possible
		if lastIndexOfNewLine == -1 {
			tmpBuffer := make([]byte, cap(buffer) * 2)
			copy(tmpBuffer, buffer)
			buffer = tmpBuffer

			log.Println("NEW BUFFER CAPACITY:", cap(buffer))
			start = cap(buffer) / 2
			continue
		}


		// for idx, piece := range strings.Split(string(buffer[:lastIndexOfNewLine + 1]), "\n") {
		// 	// log.Println(idx, piece)
		// }

		copy(buffer, buffer[lastIndexOfNewLine + 1:])
		log.Println("MODIFIED: \n",string(buffer))
		start = cap(buffer) - lastIndexOfNewLine - 1

		log.Println(start, " ", cap(buffer))
	}

	// clear buffer
	buffer = nil
	log.Println("after nil:", buffer)
}