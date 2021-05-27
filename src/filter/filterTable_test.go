package filter

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetColIndex(t *testing.T) {
	tests := []struct {
		header []string
		list   []string
		want   []int
	}{
		{[]string{"a", "b", "c", "d", "e", "f"}, []string{"e", "b", "z", "c"}, []int{4, 1, 2}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := GetColIndex(tc.header, tc.list)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
