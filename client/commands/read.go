package commands

import "log"

func Read(path string) {
	if !Exists(path) {
		log.Fatal("invalid file path " + path)
	}
}
