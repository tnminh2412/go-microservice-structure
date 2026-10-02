package random

import "testing"

func TestRandomNum(t *testing.T) {
	tests := []struct {
		name string
		no1  int64
		no2  int64
	}{
		{
			name: "test",
			no1:  1,
			no2:  10,
		},
		{
			name: "test",
			no1:  0,
			no2:  100,
		},
		{
			name: "test",
			no1:  200,
			no2:  100,
		},
		{
			name: "test",
			no1:  -200,
			no2:  -100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Between(tt.no1, tt.no2)
			t.Log(got)
		})
	}
}

func TestRandomDigit(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "test",
			length: 10,
		},
		{
			name:   "test",
			length: 6,
		},
		{
			name:   "test",
			length: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Digit(tt.length)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(got)
		})
	}
}

func BenchmarkRandomDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		got, _ := Digit(6)
		_ = got
		// b.Log(got)
	}
}
