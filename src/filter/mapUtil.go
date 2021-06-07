package filter

import "sort"

// KeyStrStr return the keys of a map map[string]string
func KeysIntPoint(mymap map[int][]Point) []int {
	keys := make([]int, 0, len(mymap))
	for k := range mymap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
