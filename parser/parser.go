package parser

import (
	"errors"
	"io"
	"strings"
	"bufio"
	"os"
	// "fmt"
	"log"
	"omise_challenges/cipher"
)


const (
	INITIAL_CAPACITY = 50
	// MAX_CAPACITY prevents from going into infinite loop
	// In case of going over MAX_CAPACITY it's expected that
	// input is not valid
	MAX_CAPACITY = 1000
)

func Parse(fileName string, outputChan chan<- string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer func () {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	reader := bufio.NewReader(file)
	rot12Reader, err := cipher.NewRot128Reader(reader)

	buffer := make([]byte, INITIAL_CAPACITY)

	index := 0
	for {
		n, err := rot12Reader.Read(buffer[index:])

		// If error is not io.EOF execution can't continue
		if err != io.EOF && err != nil {
			return err
		}

		// In case of Read doesn't populate the whole buffer
		if n != 0 && cap(buffer) - index != n {
			index += n
			continue
		}

		lastIndexOfNewLine := strings.LastIndex(string(buffer), "\n")

		// Try to keep buffer capacity as low as possible
		// When no new lines are found it's expected that buffer capacity is to small
		if lastIndexOfNewLine == -1 {
			if cap(buffer) > MAX_CAPACITY {
				return errors.New("MAX_CAPACITY Reached: check if input data has new lines")
			}

			tmpBuffer := make([]byte, cap(buffer) * 2)
			copy(tmpBuffer, buffer)
			buffer = tmpBuffer

			index = cap(buffer) / 2
			continue
		}

		// In case of more than one line is in the buffer
		for _, piece := range strings.Split(string(buffer[:lastIndexOfNewLine + 1]), "\n") {
			if piece != "" {
				outputChan <- piece
			}
		}

		// If End of File is reached parsing is done
		// Close the output channel
		if err == io.EOF {
			close(outputChan)
			break
		}

		// Rearrange memmory
		copy(buffer, buffer[lastIndexOfNewLine + 1:])
		index = cap(buffer) - lastIndexOfNewLine - 1
	}

	return nil
}
