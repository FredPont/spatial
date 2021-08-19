package ui

import (
	"fmt"
	"testing"
)

func TestInSquare(t *testing.T) {
	fmt.Println("test vulcanon In square...")
	tests := []struct {
		x, y, x1, y1, size int
		want               bool
	}{
		{0, 0, 0, 0, 10, true},
		{10, 0, 0, 0, 10, false},
		{5, 0, 0, 0, 10, true},
		{6, 0, 0, 0, 10, false},
		{5, 5, 0, 0, 10, true},
		{5, 6, 0, 0, 10, false},
		{222, 751, 218, 750, 10, true},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := inSquare(tc.x, tc.y, tc.x1, tc.y1, tc.size)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestDataColwidth(t *testing.T) {
	fmt.Println("test vulcanon In square...")
	tests := []struct {
		data [][]string
		want float32
	}{
		{[][]string{[]string{"top left very large column", "top right"}, []string{"bottom left", "bottom right"}}, 260},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := dataColwidth(tc.data)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
