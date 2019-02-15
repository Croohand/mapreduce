package fsutil

import "testing"

func TestGenerateBlockId(t *testing.T) {
	for i := 0; i < 10000; i++ {
		GenerateBlockId()
	}
}

func TestGenerateTransactionId(t *testing.T) {
	for i := 0; i < 10000; i++ {
		GenerateTransactionId()
	}
}
