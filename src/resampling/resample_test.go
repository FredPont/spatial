package resampling

import (
	"fmt"
	"log"
	"testing"
)

func TestRoundRowToInt(t *testing.T) {
	fmt.Println("test RoundRowToInt row...")
	tests := []struct {
		row     string
		indexes []int
		delim   string
		colmm   *[]ColMinMax
		want    string
	}{
		{"cellId	3.5	2.5	6.6	3.3",
			[]int{3, 4},
			"\t",
			&[]ColMinMax{{Min: 0, Max: 0}, {Min: 0, Max: 0}, {Min: 0, Max: 0}, {Min: 0, Max: 0}},
			"cellId	3.5	2.5	7	3",
		},
		{"cellId	3.5	2.5	6.6	3.3",
			[]int{1, 2, 4},
			"\t",
			&[]ColMinMax{{Min: 0, Max: 0}, {Min: 0, Max: 0}, {Min: 0, Max: 0}, {Min: 0, Max: 0}},
			"cellId	4	3	6.6	3",
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := roundRowToInt(tc.row, tc.indexes, tc.delim, tc.colmm)
			log.Println(got)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
