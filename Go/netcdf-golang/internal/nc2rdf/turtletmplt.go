package nc2rdf

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
)

type NCdataitem struct {
	Fid   int32
	Sf    float64
	Ng    float64
	Vel   float64
	Fname string
	OPID  string
}

// TurtleTemplate builds RDF in Turtle using Go's text template
func TurtleTemplate(nd NCdata) ([]byte, error) {
	fmt.Println("Generate RDF in turtle formats from UFOKN templates")

	var ba []byte

	p := prefixset()
	ba = append(ba, p...)

	for y := 0; y < len(nd.Fid); y++ {
		salt := fmt.Sprintf("%s_%s", nd.Fid[y], nd.Fname)
		h := shahash(salt)
		item := NCdataitem{Fid: nd.Fid[y], Sf: nd.Sf[y], Ng: nd.Ng[y], Vel: nd.Vel[y], Fname: nd.Fname, OPID: h}
		b, err := babyTurtle(item)
		if err != nil {
			log.Printf("baby turtle died  :(   %s", err)
		}
		ba = append(ba, b...)
	}

	return ba, nil
}

func babyTurtle(i NCdataitem) ([]byte, error) {
	tf := "./templates/nwm.rdf"

	t, err := template.New("object template").ParseFiles(tf) // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "T", i)
	if err != nil {
		log.Printf("template execution failed: %s", err)
	}

	return buf.Bytes(), err
}

func prefixset() []byte {

	p := `PREFIX schema: <http://schema.org/>
prefix dcterm: <http://purl.org/dc/terms/>
PREFIX geoparql: <http://www.opengis.net/ont/geosparql#>
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

// complete hack for the single file runs...
// generated access to the object store will not need this
func file2date(s string) string {

	return "the date pulled from the filename"
}
