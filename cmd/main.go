package main

import (
	"fmt"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	cmd := os.Args[1]

	file, err := os.Create("/Users/nathanreginato/Desktop/file.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Fprintf(file, cmd)
}
