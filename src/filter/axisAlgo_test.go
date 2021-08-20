package filter

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTicInterval(t *testing.T) {
	tests := []struct {
		yMin, yMax float64
		ticks      int
		want       []float64
	}{
		{0, 10, 5, []float64{0, 4, 8}},
		{0, 10, 4, []float64{0, 6}},
		{0, 1.0, 5, []float64{0, 0.4, 0.8}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := TicInterval(tc.yMin, tc.yMax, tc.ticks)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		yMin, yMax, y float64
		want          bool
	}{
		{0, 10, 5, true},
		{0, 10, 10, true},
		{0, 10, 0, true},
		{-10, 10, 0.5, true},
		{-10, 10, -0.5, true},
		{-10, 10, -10.5, false},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := Between(tc.yMin, tc.yMax, tc.y)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
