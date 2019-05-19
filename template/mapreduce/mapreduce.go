package mapreduce

import "hash/fnv"

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type Entry struct {
	Key, Value string
}

func Map(in <-chan string, out chan Entry) {
	for rec := range in {
		out <- Entry{"1", rec}
	}
	close(out)
}

func Partition(key string, reducers int) int {
	return int(hash(key)) % reducers
}

func Reduce(in chan Entry, out chan string) {
	for entry := range in {
		out <- entry.Key + entry.Value
	}
	close(out)
}
