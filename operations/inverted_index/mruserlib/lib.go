package mruserlib

import (
	"hash/fnv"
	"regexp"
	"strings"
)

type Entry struct {
	Key, Value string
}

func normalize(s string) string {
	r, err := regexp.Compile("[\\P{L}-_]+")
	if err != nil {
		return ""
	}
	return strings.ToLower(r.ReplaceAllString(s, ""))
}

func Map(in <-chan string, out chan<- Entry) {
	index := map[string]map[string]bool{}
	for rec := range in {
		toks := strings.SplitN(rec, "\t", 2)
		title, text := toks[0], toks[1]
		words := strings.Split(text, " ")
		for _, word := range words {
			word = normalize(word)
			if len(word) == 0 {
				continue
			}
			if _, ok := index[word]; !ok {
				index[word] = make(map[string]bool)
			}
			index[word][title] = true
		}
	}
	for word, titles := range index {
		args := []string{}
		for title := range titles {
			args = append(args, title)
		}
		joined := strings.Join(args, "\t")
		out <- Entry{word, joined}
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

func Reduce(key string, in <-chan string, out chan<- string) {
	res := map[string]bool{}
	for val := range in {
		titles := strings.Split(val, "\t")
		for _, title := range titles {
			res[title] = true
		}
	}
	args := []string{key}
	for title := range res {
		args = append(args, title)
	}
	out <- strings.Join(args, "\t")
	close(out)
}
