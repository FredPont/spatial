package ui

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestFoldChangePV(t *testing.T) {
	fmt.Println("test comparison FC and PV...")
	tests := []struct {
		table1, table2 [][]string
		colnames       []string
		want           []PVrecord
	}{
		{[][]string{{"cell1", "0.3"}, {"cell2", "0.25"}, {"cell3", "0.22"}},
			[][]string{{"cell5", "3"}, {"cell6", "2.5"}, {"cell7", "2.2"}},
			[]string{"id", "item1"},
			[]PVrecord{{"item1", 10., 0.10000000000000002, 0.10000000000000002, 3.321928094887362, 0.9999999999999999}},
		},
		{[][]string{{"cell1", "0.3"}, {"cell2", "0.25"}, {"cell3", "0.22"}},
			[][]string{{"cell5", "3.6"}, {"cell6", "3"}, {"cell7", "2.64"}},
			[]string{"id", "item1"},
			[]PVrecord{{"item1", 12., 0.10000000000000002, 0.10000000000000002, 3.584962500721156, 0.9999999999999999}},
		},
		{[][]string{{"cell1", "0.3"}, {"cell2", "0.25"}, {"cell3", "0.22"}, {"cell4", "0.1"}},
			[][]string{{"cell5", "360"}, {"cell6", "300"}, {"cell7", "264"}, {"cell8", "110"}},
			[]string{"id", "item1"},
			[]PVrecord{{"item1", 1188.5057471264367, 0.028571428571428577, 0.028571428571428577, 10.21493316424093, 1.5440680443502754}},
		},
		// {[][]string{{"cell1", "0"}, {"cell2", "0"}, {"cell3", "0"}},
		// 	[][]string{{"cell5", "3"}, {"cell6", "2.5"}, {"cell7", "2.2"}},
		// 	[]string{"id", "item1"},
		// 	[]PVrecord{},
		// },
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := foldChangePV(tc.table1, tc.table2, tc.colnames)
			log.Println(got)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestFoldChange(t *testing.T) {
	fmt.Println("test comparison FC...")
	tests := []struct {
		x1, x2 []float64
		t      bool
		want   float64
	}{
		{[]float64{0, 0, 0},
			[]float64{0, 0, 0},
			false,
			1.,
		},
		{[]float64{10, 20, 30},
			[]float64{0, 0, 0},
			true,
			1e-300,
		},
		{[]float64{0, 0, 0},
			[]float64{10, 20, 30},
			true,
			1e300,
		},
		{[]float64{5, 10, 15},
			[]float64{10, 20, 30},
			true,
			2.,
		},
		{[]float64{5, 10, 15},
			[]float64{2.5, 10, 2.5},
			true,
			0.5,
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got, err := folchange(tc.x1, tc.x2)
			log.Println(got)
			if got != tc.want || err != tc.t {
				t.Fatalf("got %v %t; want %v %t", got, err, tc.want, tc.t)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
