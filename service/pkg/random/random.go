package random

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

var r *rand.Rand

func init() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand fail")
	}

	r = rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
}

func Between(num1, num2 int64) int64 {
	if num1 == num2 {
		return num1
	}
	min := num1
	max := num2
	if num2 < num1 {
		min = num2
		max = num1
	}
	return r.Int63n(max-min+1) + min
}

var digitChars = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func Digit(length int) (string, error) {
	rep := make([]byte, length)
	cur := 0

	for cur < length {
		b := make([]byte, 2*length)
		_, err := crypto_rand.Read(b[:])
		if err != nil {
			return "", err
		}

		for _, n := range b {
			if n > 249 {
				continue
			}
			i := n % 10
			rep[cur] = digitChars[i]
			cur += 1
			if cur >= length {
				break
			}
		}
	}

	return string(rep), nil
}
