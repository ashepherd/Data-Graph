package nc2rdf

import (
	"bytes"
	"fmt"
	"log"

	"github.com/knakk/rdf"
)

// IILTriple builds a IRI, IRI, Literal triple
func IILTriple(s, p, o, c string) (string, error) {
	buf := bytes.NewBufferString("")

	newctx, err := rdf.NewIRI(c) // this should be  c
	if err != nil {
		return buf.String(), err
	}
	ctx := rdf.Context(newctx)

	sub, err := rdf.NewIRI(s)
	if err != nil {
		log.Println("Error building subject IRI")
		return buf.String(), err
	}
	pred, err := rdf.NewIRI(p)
	if err != nil {
		log.Println("Error building predicate IRI")
		return buf.String(), err
	}

	// type to xsd:double for float64
	dt, err := rdf.NewIRI("http://www.w3.org/2001/XMLSchema#double")
	if err != nil {
		log.Println("Error building predicate IRI for xsd type")
		return buf.String(), err
	}

	obj := rdf.NewTypedLiteral(o, dt) // why do typed litterials not return error?
	// if err != nil {
	// 	log.Println("Error building object literal")
	// 	return buf.String(), err
	// }

	t := rdf.Triple{Subj: sub, Pred: pred, Obj: obj}
	q := rdf.Quad{Triple: t, Ctx: ctx}

	qs := q.Serialize(rdf.NQuads)
	if s != "" && p != "" && o != "" {
		fmt.Fprintf(buf, "%s", qs)
	}
	return buf.String(), err
}

// IIITriple builds a IRI, IRI, IRI triple
func IIITriple(s, p, o, c string) (string, error) {
	buf := bytes.NewBufferString("")

	newctx, err := rdf.NewIRI(c) // this should be  c
	if err != nil {
		return buf.String(), err
	}
	ctx := rdf.Context(newctx)

	sub, err := rdf.NewIRI(s)
	if err != nil {
		log.Println("Error building subject IRI for tika triple")
		return buf.String(), err
	}
	pred, err := rdf.NewIRI(p)
	if err != nil {
		log.Println("Error building predicate IRI for tika triple")
		return buf.String(), err
	}
	obj, err := rdf.NewIRI(o)
	if err != nil {
		log.Println("Error building object literal for tika triple")
		return buf.String(), err
	}

	t := rdf.Triple{Subj: sub, Pred: pred, Obj: obj}
	q := rdf.Quad{Triple: t, Ctx: ctx}

	qs := q.Serialize(rdf.NQuads)
	if s != "" && p != "" && o != "" {
		fmt.Fprintf(buf, "%s", qs)
	}
	return buf.String(), err
}
