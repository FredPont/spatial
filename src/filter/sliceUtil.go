/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2021
*/

package filter

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

// ScaleSlice01 scale a slice between 0-1 and return the scaled 0-1 slice,  min and max
func ScaleSlice01(s []float64) ([]float64, float64, float64) {
	var norm []float64
	min, max := FindMinAndMax(s)
	for _, v := range s {
		if max != min {
			z := (v - min) / (max - min)
			norm = append(norm, z)
		} else {
			norm = append(norm, 0.)
		}
	}

	return norm, min, max
}

// ScaleSliceMinMax scale a slice between min-Max
func ScaleSliceMinMax(s []float64, min, max float64) []float64 {
	var norm []float64

	for _, v := range s {
		if max != min {
			z := (v - min) / (max - min)
			norm = append(norm, z)
		} else {
			norm = append(norm, 0.)
		}
	}

	return norm
}

// credit : https://learningprogramming.net/golang/golang-golang/find-max-and-min-of-array-in-golang/
func FindMinAndMax(a []float64) (min, max float64) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

// search str in m keys
func strInMap(str string, m map[string]bool) bool {
	_, found := m[str]
	return found
}

// StrToMap convert an array of string to map[string]bool
func StrToMap(a []string) map[string]bool {
	dic := make(map[string]bool, len(a))

	for _, x := range a {
		dic[x] = true
	}
	return dic
}

// PopIntItem return last int from []int
func PopIntItem(s []int) int {
	return s[len(s)-1]
}

// PopIntArray return the new []int witout last item
func PopIntArray(s []int) []int {
	return s[:len(s)-1]
}

// PopPointItem return last Point from []Point
func PopPointItem(s []Point) Point {
	return s[len(s)-1]
}

// StrToInt convert a string to an int
func StrToInt(s string) int {
	intVar, err := strconv.Atoi(s)
	if err != nil {
		log.Println("cannot convert string ", s, " to int !")
	}
	return intVar
}

// StrToF64 convert a string to a float64
func StrToF64(s string) float64 {
	floatVar, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("cannot convert string ", s, " to float64 !")
	}
	return floatVar
}

// ShuffleInt randomise a slice of int
func ShuffleInt(a []int) []int {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

	return a
}

// FillSliceInt fills a slice with n integers from 0 to n-1
func FillSliceInt(n int) []int {
	var slice = make([]int, n)
	for i := range slice {
		slice[i] = i
	}
	return slice
}

// DivideNB divide a number into integer parts
func DivideNB(number, numParts int) []int {

	// Calculate the quotient and remainder
	quotient := number / numParts
	remainder := number % numParts

	// Create an array to store the resulting parts
	parts := make([]int, numParts)

	// Fill the array with the quotient value
	for i := 0; i < numParts; i++ {
		parts[i] = quotient
	}

	// Distribute the remainder evenly among the parts
	for i := 0; i < remainder; i++ {
		parts[i]++
	}

	// Print the resulting parts
	//fmt.Println(parts)
	return parts
}

// DivideSlice divide a slice into  parts
func DivideSlice(slice []Point, parts int) [][]Point {

	length := len(slice)
	if length < parts {
		parts = length
	}

	result := make([][]Point, parts)
	partSize := length / parts
	remainder := length % parts

	start := 0
	for i := 0; i < parts; i++ {
		end := start + partSize
		if remainder > 0 {
			end++
			remainder--
		}

		result[i] = slice[start:end]
		start = end
	}

	return result
}

// SumSliceInt, make the sum of the n first numbers in a slice
func SumSliceInt(n int, numbers []int) int {
	// Sum of first n numbers
	sum := 0
	for i := 0; i < n; i++ {
		sum += numbers[i]
	}
	return sum
	//fmt.Printf("Sum of first %d numbers: %d", n, sum)
}
