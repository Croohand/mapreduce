package fsutil

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ValidateTransactionId(id string) bool {
	match, err := regexp.Match("[a-zA-Z0-9]{5}(\\-[a-zA-Z0-9]{5}){2}", []byte(id))
	return err == nil && match
}

func GenerateTransactionId() string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 15)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}
	id := strings.Join([]string{string(b[0:5]), string(b[5:10]), string(b[10:15])}, "-")
	if !ValidateTransactionId(id) {
		panic("generated id " + id + " not matching pattern")
	}
	return id
}
