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

// DotsInGate takes a gate and extract the scatter point that are in that gate
func DotsInGate(gate []Point, scatter map[string]Point) map[string]Point {

	cells := make(map[string]Point, 0) // cells names and XY in the gate

	for cell, point := range scatter {

		if isInside(gate, point) {
			cells[cell] = point
		}
	}
	return cells
}
