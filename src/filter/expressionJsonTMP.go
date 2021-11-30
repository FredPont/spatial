package filter

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Record store the expression variable that have to be saved in tmp files
type Record struct {
	Pts   []Point   `json:"pts"`
	Exp   []float64 `json:"Exp"`
	Min   float64   `json:"min"`
	Max   float64   `json:"Max"`
	NbPts int       `json:"nbPts"`
}

// DumpJson save expression tmp files
func DumpJson(fname string, rec Record) {
	fp, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	data, err := json.Marshal(rec)
	if err != nil {
		panic(err)
	}
	_, err = fp.Write(data)
	if err != nil {
		panic(err)
	}
}

// LoadJson read expression tmp files
func LoadJson(fname string) Record {
	var cs Record
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &cs)
	if err != nil {
		panic(err)
	}
	//fmt.Println(cs)
	return cs
}
