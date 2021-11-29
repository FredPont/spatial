package plot

import (
	"fmt"
	"reflect"
	"testing"

	"gonum.org/v1/plot/plotter"
)

func TestDensity(t *testing.T) {
	tests := []struct {
		data []float64
		n    float64
		want plotter.XYs
	}{
		{[]float64{0, 1, 2, 3, 4}, 2., plotter.XYs{plotter.XY{X: 1, Y: 2}, plotter.XY{X: 3, Y: 3}}},
		{[]float64{0, 0.5, 0.6, 1, 2, 3, 4}, 2., plotter.XYs{plotter.XY{X: 1, Y: 4}, plotter.XY{X: 3, Y: 3}}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("SumNumbers=%d", i), func(t *testing.T) {
			got := Density(tc.data, tc.n)

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
