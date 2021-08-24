package ui

import (
	"fmt"
	"lasso/src/filter"
	"reflect"
	"testing"
)

func TestFindMin(t *testing.T) {
	fmt.Println("test zoom min...")
	tests := []struct {
		imgH, imgW int
		wH, wW     float64
		want       int
	}{
		{2000, 2000, 1000, 1000, 50},
		{2000, 2000, 500, 500, 30},
		{2000, 1000, 500, 500, 50},
		{1000, 2000, 500, 500, 50},
		{1000, 2000, 300, 500, 30},
		{1000, 2000, 290, 500, 30},
		{1000, 2000, 310, 500, 40},
		{2000, 1000, 500, 310, 40},
		{2000, 2000, 1000, 1800, 90},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := findMin(tc.imgH, tc.imgW, tc.wH, tc.wW)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestZoomPolygon(t *testing.T) {
	fmt.Println("test zoom min...")
	tests := []struct {
		p    []filter.Point
		zf   float64
		want []filter.Point
	}{
		{[]filter.Point{filter.Point{800, 600}}, 2, []filter.Point{filter.Point{1600, 1200}}},
		{[]filter.Point{filter.Point{800, 600}}, 0.5, []filter.Point{filter.Point{400, 300}}},
		{[]filter.Point{filter.Point{800, 600}}, 0.12, []filter.Point{filter.Point{96, 72}}},
		{[]filter.Point{filter.Point{800, 600}}, 0.012, []filter.Point{filter.Point{10, 7}}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := filter.ZoomPolygon(tc.p, tc.zf)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
