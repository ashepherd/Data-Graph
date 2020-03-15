package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"../../internal/nc2rdf"
	"../../internal/s3fetch"
	"../../internal/urlgen"
)

func main() {
	dir, err := ioutil.TempDir("/tmp", "nwmfiles")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	start := time.Now()
	s := start.AddDate(0, 0, -4) // back three days
	e := start.AddDate(0, 0, -3) // back one day
	urls := urlgen.NameSet(s, e)

	for _, dataurl := range urls {
		fmt.Println(dataurl)

		// fp, err := fetch.GetNWM(dir, dataurl) // fetch not used for NWM, required s3 API calling
		//if err != nil {
		//	log.Println(err)
		//}
		//fmt.Println(fp)

		// convert URL into object ID we need
		u, err := url.Parse(dataurl)
		if err != nil {
			panic(err)
		}
		ps := strings.Split(u.Path, "/")
		oid := fmt.Sprintf("%s/%s/%s", ps[2], ps[3], ps[4])
		ncfn := fmt.Sprintf("%s_%s_%s", ps[2], ps[3], ps[4])
		rdffn := strings.Replace(ncfn, ".nc", ".nt", 1)

		fmt.Println("-----------------")
		fmt.Println(ncfn)
		fmt.Println(rdffn)
		fmt.Println("-----------------")

		// Fetch NetCDF file from Google S3 system hosting NWM files
		log.Printf("Fetch from Google: %s", oid)
		nwm, err := s3fetch.GetS3FP(oid)
		if err != nil {
			log.Println(err)
		}
		n, err := bytes2file(nwm, fmt.Sprintf("%s/%s", dir, ncfn))
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Fetch wrote %d bytes to %s\n", n, ncfn)
		}

		// Read NetCDF file and convert to RDF
		log.Printf("NC2RDF Reading: %s", ncfn)
		b, err := nc2rdf.ReadNC(fmt.Sprintf("%s/%s", dir, ncfn))
		if err != nil {
			log.Println(err)
		}
		n, err = bytes2file(b, fmt.Sprintf("./data/tmp/%s", rdffn))
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("NC2RDF Wrote %d bytes to %s\n", n, rdffn)
		}
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
