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
	"fmt"
	"log"
	fileutil "spatial/src/filter"
	"time"

	"github.com/akrylysov/pogreb"
)

func CSVtoPogreb() {

	t0 := time.Now()

	datafiles := fileutil.ListFiles("data")

	// remove and create db dir
	CleanDB("temp/pogreb")
	log.Println("Start ", datafiles[0], " conversion to database...")
	process(datafiles[0])

	t1 := time.Now()

	fmt.Println("")
	fmt.Println("Elapsed time : ", t1.Sub(t0))
}

func process(file string) {

	// create DB
	dbname, _ := fileutil.RemExt(file)
	db, err := pogreb.Open("temp/pogreb/"+dbname, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	// read CSV
	ReadAll(db, "data/"+file)

}
