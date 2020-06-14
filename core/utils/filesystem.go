package utils

import (
	"errors"
	"log"
	"os"
)

func OpenOrCreate(filename string) *os.File {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(errors.New("error creating output file"))
		}
		return file
	} else {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(errors.New("error opening output file"))
		}
		return file
	}
}
