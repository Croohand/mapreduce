package blockutil

import "testing"

func TestGenerateId(t *testing.T) {
	for i := 0; i < 100000; i++ {
		GenerateId()
	}
}
