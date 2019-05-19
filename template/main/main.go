package main

import (
	"log"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	if len(os.Args) < 2 {
		panic("Not enough args")
	}
	switch os.Args[1] {
	case "Map":
		doMap()
	case "Reduce":
		doReduce()
	default:
		panic("Unknown mode")
	}
}
