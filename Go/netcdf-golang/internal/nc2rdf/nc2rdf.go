package nc2rdf

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/fhs/go-netcdf/netcdf"
)

type NCdata struct {
	Fid   []int32
	Sf    []float64
	Ng    []float64
	Vel   []float64
	Fname string
}

// ReadNC reads the data in NetCDF file at filename and prints it out.
// TODO change to ReadNWM as it's not a generic nc reader.
func ReadNC(filename string) ([]byte, error) {
	fmt.Printf("Opening: %s \n", filename)

	// Open in read-only mode
	ds, err := netcdf.OpenFile(filename, netcdf.NOWRITE)
	if err != nil {
		return nil, err
	}
	defer ds.Close()

	// The following could all take place concurrently
	fid, err := getVarInt(&ds, "feature_id")
	if err != nil {
		return nil, err
	}

	sf, err := getVarFloat(&ds, "streamflow", 0.01)
	if err != nil {
		return nil, err
	}

	ng, err := getVarFloat(&ds, "nudge", 0.01)
	if err != nil {
		return nil, err
	}

	vel, err := getVarFloat(&ds, "velocity", 0.01)
	if err != nil {
		return nil, err
	}

	// form a context (named graph) from the filename
	base := filepath.Base(filename)
	// rgxbase := strings.ReplaceAll(base, ".", "_")
	// gctx := fmt.Sprintf("https://ufokn.org/ctx/nwm/%s", rgxbase) // define a quad context string (will be an IRI)

	nd := NCdata{Fid: fid, Sf: sf, Ng: ng, Vel: vel, Fname: base}

	// rb, err := TripleGen(nd, gctx)
	// log.Println(len(rb))

	rb2, err := TurtleTemplate(nd)
	log.Println(len(rb2))

	return rb2, err
}
