package blockutil

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type PathInfo []BlockInfo

type BlockInfo struct {
	Id           string
	Lower, Upper int
	Slaves       []string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ValidateId(id string) bool {
	match, err := regexp.Match("[a-zA-Z0-9]{4}(\\-[a-zA-Z0-9]{4}){3}", []byte(id))
	return err == nil && match
}

func ValidateFilePath(path string) bool {
	return len(path) > 0
}

func GenerateId() string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}
	id := strings.Join([]string{string(b[0:4]), string(b[4:8]), string(b[8:12]), string(b[12:16])}, "-")
	if !ValidateId(id) {
		panic("generated id " + id + " not matching pattern")
	}
	return id
}
