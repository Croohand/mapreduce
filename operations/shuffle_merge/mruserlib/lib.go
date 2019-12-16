package mruserlib

import (
	"math/rand"
	"time"
)

type Entry struct {
	Key, Value string
}

func Map(in <-chan string, out chan<- Entry) {
	rand.Seed(time.Now().UTC().UnixNano())
	entries := []Entry{}
	for rec := range in {
		entries = append(entries, Entry{"0", rec})
	}
	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})
	for _, rec := range entries {
		out <- rec
	}
	close(out)
}

func Partition(key string, reducers int) int {
	return rand.Intn(reducers)
}

func Reduce(key string, in <-chan string, out chan<- string) {
	for val := range in {
		out <- val
	}
	close(out)
}
