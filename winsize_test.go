package main

import (
	"fmt"
	"testing"
)

func TestSetMinWindow(t *testing.T) {
	fmt.Println("test setMinWindow...")
	tests := []struct {
		prefSize float64
		imgSize  int
		want     float32
	}{
		{0., 2000, 500.},
		{600., 2000, 600.},
		{800., 700, 700.},
		{800., 400, 400.},
		{11., 400, 50.},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := setMinWindow(tc.prefSize, tc.imgSize)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
