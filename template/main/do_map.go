package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"mruserlib"
)

func writeMapInput(inFile string, in chan<- string, done chan<- bool) {
	file, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if len(line) > 0 {
			in <- line
		}
	}
	close(in)
	if sc.Err() != nil {
		panic(sc.Err())
	}
	done <- true
}

func getOutFile(outputDir string, reducer int) string {
	return filepath.Join(outputDir, strconv.Itoa(reducer))
}

func fetchMapOutput(outputDir string, reducers int, out <-chan mruserlib.Entry, done chan<- bool) {
	files := map[int]*os.File{}
	for r := 0; r < reducers; r++ {
		file, err := os.Create(getOutFile(outputDir, r))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		files[r] = file
	}
	for entry := range out {
		if len(entry.Key) == 0 {
			panic("Error empty key")
		}
		if len(entry.Value) == 0 {
			panic("Error empty value")
		}
		r := mruserlib.Partition(entry.Key, reducers)
		if r < 0 || r >= reducers {
			panic("Bad partition function output")
		}
		_, err := files[r].Write([]byte(entry.Key + "\t" + entry.Value + "\n"))
		if err != nil {
			panic(err)
		}
	}
	done <- true
}

func sortFiles(outputDir string, reducers int) {
	for r := 0; r < reducers; r++ {
		file := getOutFile(outputDir, r)
		cmd := exec.Command("sort", file, "-o", file)
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

func doMap() {
	if len(os.Args) != 5 {
		panic("Invalid arguments for Map")
	}
	inputFile := os.Args[2]
	outputDir := os.Args[3]
	reducers, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic(err)
	}
	if reducers <= 0 {
		panic("Invalid number of reducers")
	}
	in := make(chan string)
	out := make(chan mruserlib.Entry)
	done := make(chan bool)
	go writeMapInput(inputFile, in, done)
	go mruserlib.Map(in, out)
	go fetchMapOutput(outputDir, reducers, out, done)
	<-done
	<-done
	sortFiles(outputDir, reducers)
}
