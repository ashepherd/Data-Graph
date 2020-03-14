package nc2rdf

import (
	"bytes"
	"fmt"
	"log"

	"github.com/fhs/go-netcdf/netcdf"
)

// ReadNC reads the data in NetCDF file at filename and prints it out.
// TODO change to ReadNWM as it's not a generic nc reader.
func ReadNC(filename string) ([]byte, error) {
	// Open read-only mode, why not...
	ds, err := netcdf.OpenFile(filename, netcdf.NOWRITE)
	if err != nil {
		return nil, err
	}
	defer ds.Close()

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

	buf := bytes.NewBufferString("")

	// Print out the data
	for y := 0; y < len(fid); y++ {
		// format an ID for this resource
		rid := fmt.Sprintf("https://ufokn.org.x/id/nwm/%d", fid[y])
		// define a quad context string (will be an IRI)
		gctx := "https://ufokn.org/ctx/nwm"

		b1, err := IILTriple(rid, "https://ufokn.org/voc/1/streamflow", fmt.Sprintf("%f", sf[y]), gctx)
		if err != nil {
			log.Println(err)
		}
		b2, err := IILTriple(rid, "https://ufokn.org/voc/1/nudge", fmt.Sprintf("%f", ng[y]), gctx)
		if err != nil {
			log.Println(err)
		}
		b3, err := IILTriple(rid, "https://ufokn.org/voc/1/velocity", fmt.Sprintf("%f", vel[y]), gctx)
		if err != nil {
			log.Println(err)
		}
		b4, err := IIITriple(rid, "http://rdfs/type", "https://ufokn.org/voc/1/NWM", gctx)
		if err != nil {
			log.Println(err)
		}

		fmt.Fprintf(buf, "%s", b1)
		fmt.Fprintf(buf, "%s", b2)
		fmt.Fprintf(buf, "%s", b3)
		fmt.Fprintf(buf, "%s", b4)
	}

	// fmt.Print(buf.String())
	return buf.Bytes(), err
}
