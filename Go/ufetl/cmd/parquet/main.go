package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"reflect"
	"text/template"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

// <pyarrow._parquet.ParquetSchema object at 0x7ffa0a003128>
type ufokn struct {
	ID        string  `parquet:"name=id, type=UTF8, encoding=PLAIN_DICTIONARY"`        // BYTE_ARRAY String
	Type      string  `parquet:"name=type, type=UTF8, encoding=PLAIN_DICTIONARY"`      // BYTE_ARRAY String
	Name      string  `parquet:"name=name, type=UTF8, encoding=PLAIN_DICTIONARY"`      // BYTE_ARRAY String
	Hand      float64 `parquet:"name=hand, type=UTF8, encoding=PLAIN_DICTIONARY"`      // DOUBLE
	Offset    float64 `parquet:"name=offset, type=UTF8, encoding=PLAIN_DICTIONARY"`    // DOUBLE
	Featureid string  `parquet:"name=featureid, type=UTF8, encoding=PLAIN_DICTIONARY"` // DOUBLE
	X         float64 `parquet:"name=x, type=DOUBLE"`                                  // DOUBLE
	Y         float64 `parquet:"name=y, type=DOUBLE"`
	Z         float64 `parquet:"name=z, type=DOUBLE"`
	A         float64 `parquet:"name=a, type=DOUBLE"`
	B         float64 `parquet:"name=b, type=DOUBLE"`
	Source    string  `parquet:"name=source, type=UTF8, encoding=PLAIN_DICTIONARY"` // BYTE_ARRAY String
	CKey      string  `parquet:"name=ckey, type=UTF8, encoding=PLAIN_DICTIONARY"`   // BYTE_ARRAY String

	// DOUBLE
}

var floatType = reflect.TypeOf(float64(0))

func main() {
	//	var data []*ufokn
	log.Println("Reading file")

	b, err := test()
	if err != nil {
		log.Println(err)
	}

	out := "./data/output/parquet.rdf"
	n, err := bytes2file(b, out)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("NC2RDF Wrote %d bytes to %s\n", n, out)
	}

	//all, err := myreadPartialParquet(50835, 0)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//log.Println("Printing all data")
	//for _, a := range all {
	//	fmt.Println(a)
	//}
}

func test() ([]byte, error) {
	fr, err := local.NewLocalFileReader("./data/input/jefferson_complete_risk_pts_test.parquet")
	if err != nil {
		log.Println(err)
	}

	// Get a column reader and loop over the rows.
	pr, err := reader.NewParquetColumnReader(fr, 4) // 4 is the number parallel which I suspect is a go routine count?
	if err != nil {
		log.Println(err)
	}

	num := pr.GetNumRows()
	d0, _, _, err := pr.ReadColumnByIndex(0, num)
	d1, _, _, err := pr.ReadColumnByIndex(1, num)
	d2, _, _, err := pr.ReadColumnByIndex(2, num)
	d3, _, _, err := pr.ReadColumnByIndex(3, num)
	d4, _, _, err := pr.ReadColumnByIndex(4, num)
	d5, _, _, err := pr.ReadColumnByIndex(5, num)
	d6, _, _, err := pr.ReadColumnByIndex(6, num)
	d7, _, _, err := pr.ReadColumnByIndex(7, num)
	if err != nil {
		log.Println(err)
	}

	// Search testing to get A, B and Hand for a featureID
	//for k, _ := range d5 {
	//	// if strings.Compare(v, "1114365.0") {
	//	if d5[k] == 1114365.0 {
	//		fmt.Print(k)
	//		fmt.Println(d5[k])
	//	}
	//}

	// Set up templates
	tf := "./templates/parquet.rdf"
	t, err := template.New("object template").ParseFiles(tf) // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	// set up byte array
	var ba []byte

	// load RDF prefix lines to array
	p := prefixset()
	ba = append(ba, p...)

	for i := range d2 {
		xf, _ := getFloat(d6[i])
		yf, _ := getFloat(d7[i])
		handf, _ := getFloat(d3[i])
		offsetf, _ := getFloat(d4[i])

		us := ufokn{ID: fmt.Sprintf("%v", d0[i]),
			Type:      fmt.Sprintf("%v", d1[i]),
			Name:      fmt.Sprintf("%v", d2[i]),
			Hand:      handf,
			Offset:    offsetf,
			Featureid: fmt.Sprintf("%.0f", d5[i]),
			X:         xf,
			Y:         yf,
			A:         0.28,
			B:         0.38,
			Source:    "OSM",
			CKey:      shahash(fmt.Sprintf("%s%s%s", "OSM", d0[i], d5[i]))}

		// fmt.Println(us)

		var buf bytes.Buffer
		err = t.ExecuteTemplate(&buf, "T", us)
		if err != nil {
			log.Printf("template execution failed: %s", err)
		}

		bb := buf.Bytes()
		ba = append(ba, bb...)
	}

	fr.Close()
	return ba, nil
}

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func prefixset() []byte {

	p := `PREFIX schema: <http://schema.org/>
prefix dcterm: <http://purl.org/dc/terms/>
PREFIX geo: <http://www.w3.org/2003/01/geo/wgs84_pos#>
PREFIX geosparql: <http://www.opengis.net/ont/geosparql#>
PREFIX geofunc: <http://www.opengis.net/def/function/geosparql/>
PREFIX osm: <http://schema.ufokn.org/osm/v1/>
PREFIX owl: <http://www.w3.org/2002/07/owl#>
PREFIX nwm: <https://schema.ufokn.org/nvm/v1/>
PREFIX prov: <http://www.w3.org/ns/prov#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns##>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX sf: <http://www.opengis.net/ont/sf#>
PREFIX ufokn: <http://schema.ufokn.org/core/v1/>
PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

	`
	return []byte(p)
}

// sha (or md5) hash for a string for opaque (esk) IDs
func shahash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
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

// --- below here is not used....

func test2() {
	fr, err := local.NewLocalFileReader("jefferson_complete_risk_pts_test.parquet")
	if err != nil {
		log.Println(err)
	}

	pr0, err := reader.NewParquetReader(fr, new(ufokn), 4)
	if err != nil {
		log.Println(err)
	}

	log.Println(pr0.GetNumRows())
	u := make([]*ufokn, pr0.GetNumRows()*24)

	if err = pr0.Read(&u); err != nil {
		log.Println(" -------------------   check mark 0.5-----------------------")
		log.Println(err)
	}

	pr0.ReadStop()

	for i := range u {
		fmt.Print(u[i])
	}
	fr.Close()
}

func myreadPartialParquet(pageSize, page int) ([]*ufokn, error) {
	fr, err := local.NewLocalFileReader("jefferson_complete_risk_pts_test.parquet")
	if err != nil {
		return nil, err
	}
	pr, err := reader.NewParquetReader(fr, new(ufokn), int64(pageSize))
	if err != nil {
		return nil, err
	}
	pr.SkipRows(int64(pageSize * page))
	u := make([]*ufokn, pageSize)
	if err = pr.Read(&u); err != nil {
		return nil, err
	}
	pr.ReadStop()
	fr.Close()
	return u, nil
}

func myreadParquet(recordNumber int64) ([]*ufokn, error) {
	fr, err := local.NewLocalFileReader("jefferson_complete_risk_pts_test.parquet")
	if err != nil {
		return nil, err
	}
	pr, err := reader.NewParquetReader(fr, new(ufokn), recordNumber)
	if err != nil {
		return nil, err
	}
	u := make([]*ufokn, recordNumber)
	if err = pr.Read(&u); err != nil {
		return nil, err
	}
	pr.ReadStop()
	fr.Close()
	return u, nil
}
