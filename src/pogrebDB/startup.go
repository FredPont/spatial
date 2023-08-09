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
	"spatial/src/filter"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func InitPogreb() {
	// get the user preference for using the database
	pref := fyne.CurrentApp().Preferences()
	useDBpref := binding.BindPreferenceBool("useDataBase", pref)
	useDB, err := useDBpref.Get()
	check(err)
	//log.Println("use DB ", useDB)
	// if the user want to use a database the database is created if not exist
	if useDB {
		dataFiles := filter.ListFiles("data/")
		dataInUse, _ := filter.RemExt(dataFiles[0])
		DBFiles := filter.ListFiles("temp/pogreb/")
		DBfile := ""
		if len(DBFiles) != 0 {
			DBfile = DBFiles[0]

			//log.Println("data ", dataInUse, " DB ", DBfile)
			if DBfile == dataInUse {
				fmt.Println("A database with the name ", DBfile, " already exist. Do you want to replace it (y/n) ? [n]")
				var input string
				fmt.Scanf("%s", &input)
				if input == "y" || input == "Y" {
					CSVtoPogreb() // convert current csv file to pogreb DB
				}
				return
			}

		} else {
			CSVtoPogreb() // if there is no pogreb file, convert current csv file to pogreb DB
		}

	}
	return

}
