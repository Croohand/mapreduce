package mruserlib

import (
	"hash/fnv"
	"strconv"
	"strings"
)

type Entry struct {
	Key, Value string
}

func Map(in <-chan string, out chan<- Entry) {
	wordsCount := map[string]int{}
	for rec := range in {
		words := strings.Split(rec, " ")
		for _, word := range words {
			if len(word) > 0 {
				wordsCount[word]++
			}
		}
	}
	for word, count := range wordsCount {
		out <- Entry{word, strconv.Itoa(count)}
	}
	close(out)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func Partition(key string, reducers int) int {
	return int(hash(key)) % reducers
}

func Reduce(in <-chan Entry, out chan<- string) {
	wordsCount := map[string]int{}
	for entry := range in {
		count, err := strconv.Atoi(entry.Value)
		if err != nil {
			continue
		}
		wordsCount[entry.Key] += count
	}
	for word, count := range wordsCount {
		out <- word + ": " + strconv.Itoa(count)
	}
	close(out)
}
