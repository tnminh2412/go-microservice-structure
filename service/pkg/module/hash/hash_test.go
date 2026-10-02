package hash

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	n := 100
	m := 10

	fakeNames := make([]string, n)
	for i := 0; i < n; i++ {
		fakeNames[i] = gofakeit.Username()
	}

	for _, name := range fakeNames {
		hashNum := Hash(name)
		shardNo := hashNum % uint32(m)
		assert.Less(t, shardNo, uint32(m))
	}
}
