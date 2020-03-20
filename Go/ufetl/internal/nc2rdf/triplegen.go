package nc2rdf

import (
	"bytes"
	"fmt"
	"log"

	"github.com/knakk/rdf"
)

// TripleGen gernates triples from a netcdf file contents for the National Water Model
func TripleGen(nd NCdata, gctx string) ([]byte, error) {
	buf := bytes.NewBufferString("")
	for y := 0; y < len(nd.Fid); y++ {
		rid := fmt.Sprintf("https://ufokn.org.x/id/nwm/%d", nd.Fid[y]) // format an ID for this resource

		b1, err := IILTriple(rid, "http://schema.ufokn.org/core/v1/streamflow", fmt.Sprintf("%f", nd.Sf[y]), gctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		b2, err := IILTriple(rid, "http://schema.ufokn.org/core/v1/nudge", fmt.Sprintf("%f", nd.Ng[y]), gctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		b3, err := IILTriple(rid, "http://schema.ufokn.org/core/v1/velocity", fmt.Sprintf("%f", nd.Vel[y]), gctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		b4, err := IIITriple(rid, "http://rdfs/type", "http://schema.ufokn.org/core/v1/NWM", gctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		fmt.Fprintf(buf, "%s", b1)
		fmt.Fprintf(buf, "%s", b2)
		fmt.Fprintf(buf, "%s", b3)
		fmt.Fprintf(buf, "%s", b4)
	}

	// fmt.Print(buf.String())
	return buf.Bytes(), nil
}

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

// IIBTriple builds a IRI, IRI, Blank triple
func IIBTriple(s, p, o, c string) (string, error) {
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
	obj, err := rdf.NewBlank(o)
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

// BILTriple builds a Blank, IRI, Literal triple
func BILTriple(s, p, o, c string) (string, error) {
	buf := bytes.NewBufferString("")

	newctx, err := rdf.NewIRI(c) // this should be  c
	if err != nil {
		return buf.String(), err
	}
	ctx := rdf.Context(newctx)

	sub, err := rdf.NewBlank(s)
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

// BIITriple builds a Blank, IRI, IRI triple
func BIITriple(s, p, o, c string) (string, error) {
	buf := bytes.NewBufferString("")

	newctx, err := rdf.NewIRI(c) // this should be  c
	if err != nil {
		return buf.String(), err
	}
	ctx := rdf.Context(newctx)

	sub, err := rdf.NewBlank(s)
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
