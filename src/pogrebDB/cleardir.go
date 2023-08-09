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
 (c) Frederic Pont 2022
*/

package pogrebDB

import (
	"log"
	"os"
)

// CleanDB remove and create a directory
func CleanDB(dir string) {
	ClearDB(dir)
	MkDir(dir)
}

// remove DB dir
func ClearDB(dir string) {
	// clear db directory
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
}

// create DB dir
func MkDir(dir string) {
	//Create a folder/directory at a full qualified path
	err := os.Mkdir(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
