package filter

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDivideNB(t *testing.T) {
	tests := []struct {
		number, numParts int
		want             []int
	}{
		{10, 3, []int{4, 3, 3}},
		{219291, 4, []int{54823, 54823, 54823, 54822}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("SumNumbers=%d", i), func(t *testing.T) {
			got := DivideNB(tc.number, tc.numParts)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}

}

func TestSumSliceInt(t *testing.T) {
	tests := []struct {
		n      int
		number []int
		want   int
	}{
		{2, []int{4, 3, 3}, 7},
		{0, []int{4, 3, 3}, 0},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("SumNumbers=%d", i), func(t *testing.T) {
			got := SumSliceInt(tc.n, tc.number)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}

}
