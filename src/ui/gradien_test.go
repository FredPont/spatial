package ui

import (
	"fmt"
	"testing"
)

func TestWRgradien(t *testing.T) {
	fmt.Println("test red white gradien...")
	tests := []struct {
		value float64
		want  RGB
	}{
		{0., RGB{255, 245, 240}},
		{1., RGB{103, 0, 13}},
		{0.5, RGB{250, 105, 75}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := WRgradien(tc.value)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestTurboGradien(t *testing.T) {
	fmt.Println("test red white gradien...")
	tests := []struct {
		value float64
		want  RGB
	}{
		{0., RGB{35, 23, 27}},
		{1., RGB{144, 12, 0}},
		{0.5, RGB{149, 251, 81}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := TurboGradien(tc.value)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
