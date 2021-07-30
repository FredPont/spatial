package filter

import (
	"lasso/src/filter/stats"
	"log"
)

// pvalue for MannWhitney test
func PvMannWhitney(x1, x2 []float64) (float64, bool) {
	s, err := stats.MannWhitneyUTest(x1, x2, 0)
	if err != nil {
		log.Println("MannWhitney cannot be calculated", err, "group1 ", x1, "group2 ", x2)
		return 1., false
	}
	//fmt.Println(s.U, s.P, err)
	return s.P, true
}
