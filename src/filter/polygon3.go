// algorithm from https://stackoverflow.com/questions/217578/how-can-i-determine-whether-a-2d-point-is-within-a-polygon?page=2&tab=votes#tab-top
// https://wrf.ecse.rpi.edu/Research/Short_Notes/pnpoly.html
// W. Randolph Franklin (WRF)
/*License to Use
Copyright (c) 1970-2003, Wm. Randolph Franklin

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimers.
Redistributions in binary form must reproduce the above copyright notice in the documentation and/or other materials provided with the distribution.
The name of W. Randolph Franklin may not be used to endorse or promote products derived from this Software without specific prior written permission.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package filter

// Point = XY coordinates of a point
type Point struct {
	X, Y int
}

// test if testp is inside the polygon
func isInside(polygon []Point, testp Point) bool {
	if len(polygon) < 1 {
		return false
	}
	minX := polygon[0].X
	maxX := polygon[0].X
	minY := polygon[0].Y
	maxY := polygon[0].Y

	for _, p := range polygon {
		minX = Min(p.X, minX)
		maxX = Max(p.X, maxX)
		minY = Min(p.Y, minY)
		maxY = Max(p.Y, maxY)
	}

	if testp.X < minX || testp.X > maxX || testp.Y < minY || testp.Y > maxY {
		return false
	}

	inside := false
	j := len(polygon) - 1
	for i := 0; i < len(polygon); i++ {
		if (polygon[i].Y > testp.Y) != (polygon[j].Y > testp.Y) && testp.X < (polygon[j].X-polygon[i].X)*(testp.Y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X {
			inside = !inside
		}
		j = i
	}

	return inside
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
