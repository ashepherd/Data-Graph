package main

import (
	"log"
	"os"

	// "os"

	"../../internal/nc2rdf"
)

// Some of the files are in an S3 API based system at Google.  We can get by URL but
// it might be easier to use a direct object link later as we can avoid
// the local marshalling and deal with the objects as byte streams

func main() {

	// dir, err := ioutil.TempDir("/tmp", "nwmfiles")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer os.RemoveAll(dir)

	// start := time.Now()
	// s := start.AddDate(0, 0, -4) // back three days
	// e := start.AddDate(0, 0, -2) // back one day
	// urls := urlgen.NameSet(s, e)
	// for _, dataurl := range urls {
	// 	// read the URL to scratch space
	// 	// process said URL, now file
	// 	fp, err := fetch.GetNWM(dir, dataurl)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println(dataurl)
	// 	fmt.Println(fp)
	// 	b, err := nc2rdf.ReadNC("./data/input/test.nc")
	// 	if err != nil {
	// 		log.Fatalf("reading example file failed: %v\n", err)
	// 	}
	// 	fmt.Println(len(b))

	// }

	// RDF format testing snippet
	log.Println("Reading ./data/input/test.nc")
	b, err := nc2rdf.ReadNC("./data/input/test.nc")
	if err != nil {
		log.Println(err)
	}
	n, err := string2file(b, "./data/output/test.nq")
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Wrote %d bytes\n", n)
	}

}

func string2file(b []byte, fn string) (int, error) {
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		return 0, err
	}

	n, err := f.Write(b)
	return n, err
}
