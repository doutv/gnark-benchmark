package utils

import (
	"io"
	"os"

	gnark_io "github.com/consensys/gnark/io"
)

func WriteToFile(data io.WriterTo, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = data.WriteTo(file)
	if err != nil {
		panic(err)
	}
}

func ReadFromFile(data io.ReaderFrom, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Use the ReadFrom method to read the file's content into data.
	if _, err := data.ReadFrom(file); err != nil {
		panic(err)
	}
}

// faster than readFromFile
func UnsafeReadFromFile(data gnark_io.UnsafeReaderFrom, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := data.UnsafeReadFrom(file); err != nil {
		panic(err)
	}
}
