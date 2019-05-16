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

type TransactionType int

const (
	TxTypeRead TransactionType = iota
	TxTypeWrite
)

func (txType TransactionType) String() string {
	switch txType {
	case TxTypeRead:
		return "Read"
	case TxTypeWrite:
		return "Write"
	default:
		return "Unknown"
	}
}

type Transaction struct {
	Paths      []string
	LastUpdate time.Time
	TxType     TransactionType
}

type TransactionEx struct {
	Id string
	Transaction
}

type TxIds map[string]struct{}

const transactionWaitTime = 60

func NewTransaction(paths []string, txType TransactionType) *Transaction {
	return &Transaction{paths, time.Now(), txType}
}

func (tx *Transaction) IsAlive() bool {
	return time.Since(tx.LastUpdate).Seconds() <= transactionWaitTime
}

func (tx *Transaction) Update() {
	tx.LastUpdate = time.Now()
}

func ValidateTransactionId(id string) bool {
	if id == "" {
		return true
	}
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
		panic("Generated id " + id + " not matching pattern")
	}
	return id
}
