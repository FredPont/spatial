package filter

import "sort"

// KeysIntPoint return the keys of a map map[string]string
func KeysIntPoint(mymap map[int][]Point) []int {
	keys := make([]int, 0, len(mymap))
	for k := range mymap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// PopPoints remove last element of slice [][]Point = remove last gate
func PopPoints(slice [][]Point) [][]Point {
	if len(slice) > 0 {
		slice = slice[:len(slice)-1]
	}
	return slice
}
