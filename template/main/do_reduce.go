package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"mruserlib"
)

type chunk struct {
	scanner  *bufio.Scanner
	curEntry mruserlib.Entry
}

func (ch *chunk) nextEntry() bool {
	if !ch.scanner.Scan() {
		return false
	}
	line := ch.scanner.Text()
	toks := strings.SplitN(line, "\t", 2)
	if len(toks) != 2 {
		panic("Bad entry " + line)
	}
	ch.curEntry = mruserlib.Entry{toks[0], toks[1]}
	return true
}

type PriorityQueue []*chunk

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].curEntry.Key < pq[j].curEntry.Key
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*chunk))
}

func (pq *PriorityQueue) Pop() interface{} {
	n := pq.Len()
	x := (*pq)[n-1]
	*pq = (*pq)[0 : n-1]
	return x
}

func fetchReduceOutput(out <-chan string, done chan<- bool) {
	for rec := range out {
		fmt.Println(rec)
	}
	done <- true
}

func doReduce() {
	if len(os.Args) != 3 {
		panic("Invalid arguments for Reduce")
	}
	inputDir := os.Args[2]
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		return
	}
	pq := make(PriorityQueue, 0)
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}
		filePath := filepath.Join(inputDir, fileInfo.Name())
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		ch := &chunk{scanner: bufio.NewScanner(file)}
		if ch.nextEntry() {
			heap.Push(&pq, ch)
		}
	}
	prevKey := ""
	var in chan mruserlib.Entry
	var out chan string
	var done chan bool
	for pq.Len() > 0 {
		ch := heap.Pop(&pq).(*chunk)
		if ch.curEntry.Key != prevKey {
			if prevKey != "" {
				close(in)
				<-done
			}
			prevKey = ch.curEntry.Key
			in = make(chan mruserlib.Entry)
			out = make(chan string)
			done = make(chan bool)
			go mruserlib.Reduce(in, out)
			go fetchReduceOutput(out, done)
		}
		in <- ch.curEntry
		if ch.nextEntry() {
			heap.Push(&pq, ch)
		}
	}
	if prevKey != "" {
		close(in)
		<-done
	}
}
