package mask

import "testing"

func TestUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		no2      int64
	}{
		{
			name:     "test0",
			username: "abc",
		},
		{
			name:     "test1",
			username: "playermr",
		},
		{
			name:     "test2",
			username: "multimmm",
		},
		{
			name:     "test3",
			username: "lonngngngnngngngn",
		},
		{
			name:     "test4",
			username: "123456789",
		},
		{
			name:     "test5",
			username: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Username(tt.username)
			t.Log(got)
		})
	}
}
