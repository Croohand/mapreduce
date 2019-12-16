package mruserlib

import (
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
)

type Entry struct {
	Key, Value string
}

func normalize(s string) string {
	r, err := regexp.Compile("[\t,. !#?()<>]+")
	if err != nil {
		return ""
	}
	return strings.ToLower(r.ReplaceAllString(s, ""))
}

func Map(in <-chan string, out chan<- Entry) {
	wordsCount := map[string]int{}
	for rec := range in {
		words := strings.Split(rec, " ")
		for _, word := range words {
			w := normalize(word)
			if len(w) > 0 {
				wordsCount[w]++
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

func prependZeros(s string) string {
	return strings.Repeat("0", 9-len(s)) + s
}

func Reduce(key string, in <-chan string, out chan<- string) {
	res := 0
	for val := range in {
		count, err := strconv.Atoi(val)
		if err != nil {
			continue
		}
		res += count
	}
	out <- prependZeros(strconv.Itoa(res)) + ": " + key
	close(out)
}
