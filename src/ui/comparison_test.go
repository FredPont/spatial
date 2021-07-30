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
			[]PVrecord{{"item1", 10., 0.10000000000000002, 0.10000000000000002}},
		},
		{[][]string{{"cell1", "0.3"}, {"cell2", "0.25"}, {"cell3", "0.22"}},
			[][]string{{"cell5", "3.6"}, {"cell6", "3"}, {"cell7", "2.64"}},
			[]string{"id", "item1"},
			[]PVrecord{{"item1", 12., 0.10000000000000002, 0.10000000000000002}},
		},
		{[][]string{{"cell1", "0.3"}, {"cell2", "0.25"}, {"cell3", "0.22"}, {"cell4", "0.1"}},
			[][]string{{"cell5", "360"}, {"cell6", "300"}, {"cell7", "264"}, {"cell8", "110"}},
			[]string{"id", "item1"},
			[]PVrecord{{"item1", 1188.5057471264367, 0.028571428571428577, 0.028571428571428577}},
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
