package main

import (
	"log"
	"os"

	"../../internal/nc2rdf"
)

func main() {
	// TODO

	// TODO  open all files and return pointers to them

	// TODO read all files to a struct

	// TODO value add processess

	// TODO serialize

	// The Google hosted files use an S3 API that requires credentials but there is no cost.
	// You need to generate a credentials JSON file and note it in your env.  Though later I will use
	// the static file location approach.
	//googleSource()

	// The local source simply reads a NetCDF NWM based file and converts it.
	localSource("./data/input/nwm.20200210_analysis_assim_nwm.t00z.analysis_assim.channel_rt.tm00.conus.nc", "./data/output/nwm_test.rdf")
}

func localSource(in, out string) {
	// Read NetCDF file and convert to RDF
	log.Printf("NC2RDF Reading: %s", in)
	b, err := nc2rdf.ReadNC(in)
	if err != nil {
		log.Println(err)
	}
	n, err := bytes2file(b, out)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("NC2RDF Wrote %d bytes to %s\n", n, out)
	}
}

func bytes2file(b []byte, fn string) (int, error) {
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		return 0, err
	}

	n, err := f.Write(b)
	f.Close()
	return n, err
}
