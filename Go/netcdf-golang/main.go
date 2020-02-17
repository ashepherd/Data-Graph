package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/fhs/go-netcdf/netcdf"
	"github.com/knakk/rdf"
)

func main() {
	// Get the files...   (will be S3 API based)
	//fileset()

	// Create example file
	//	filename := "gopher.nc"
	//	if err := CreateExampleFile(filename); err != nil {
	//		log.Fatalf("creating example file failed: %v\n", err)
	//	}

	// Open and read example file
	//	if err := ReadExampleFile(filename); err != nil {
	//		log.Fatalf("reading example file failed: %v\n", err)
	//	}

	if err := ReadTest("test.nc"); err != nil {
		log.Fatalf("reading example file failed: %v\n", err)
	}

}

// ReadTest reads the data in NetCDF file at filename and prints it out.
func ReadTest(filename string) error {
	// Open example file in read-only mode. The dataset is returned.
	ds, err := netcdf.OpenFile(filename, netcdf.NOWRITE)
	if err != nil {
		fmt.Println("open file")
		return err
	}
	defer ds.Close()

	fid, err := getVarInt(&ds, "feature_id")
	if err != nil {
		fmt.Println(err)
		return err
	}

	sf, err := getVarFloat(&ds, "streamflow", 0.01)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ng, err := getVarFloat(&ds, "nudge", 0.01)
	if err != nil {
		fmt.Println(err)
		return err
	}

	vel, err := getVarFloat(&ds, "velocity", 0.01)
	if err != nil {
		fmt.Println(err)
		return err
	}

	buf := bytes.NewBufferString("")

	// Print out the data
	for y := 0; y < len(fid); y++ {

		// fmt.Printf(" --  %d %f \n", fid[y], sf[y])
		b, err := IILTriple(fmt.Sprintf("https://foo.x/%d", fid[y]), "https://ufokn.org/voc/streamflow", fmt.Sprintf("%f", sf[y]), "https://ufokn.org/contextquad")
		if err != nil {
			log.Println(err)
		}
		b2, err := IILTriple(fmt.Sprintf("https://foo.x/%d", fid[y]), "https://ufokn.org/voc/nudge", fmt.Sprintf("%f", ng[y]), "https://ufokn.org/contextquad")
		if err != nil {
			log.Println(err)
		}
		b3, err := IILTriple(fmt.Sprintf("https://foo.x/%d", fid[y]), "https://ufokn.org/voc/velocity", fmt.Sprintf("%f", vel[y]), "https://ufokn.org/contextquad")
		if err != nil {
			log.Println(err)
		}
		// TODO add in the IIITriple function for the type URI
		b4, err := IILTriple(fmt.Sprintf("https://foo.x/%d", fid[y]), "http://rdfs/type", "https://ufokn.org/voc/NETCDFType", "https://ufokn.org/contextquad")
		if err != nil {
			log.Println(err)
		}

		fmt.Fprintf(buf, "%s", b)
		fmt.Fprintf(buf, "%s", b2)
		fmt.Fprintf(buf, "%s", b3)
		fmt.Fprintf(buf, "%s", b4)

	}

	fmt.Print(buf.String())

	return nil

}

func fileset() {
	start := time.Now()
	s := start.AddDate(0, 0, -3)
	e := start.AddDate(0, 0, -1)

	for rd := rangeDate(s, e); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		d := fmt.Sprint(date.Format("20060102"))
		for i := 0; i <= 23; i++ {
			t := fmt.Sprintf("%02d", i)
			fmt.Printf("https://storage.cloud.google.com/national-water-model/nwm.%s/analysis_assim/nwm.t%sz.analysis_assim.channel_rt.tm00.conus.nc\n", d, t)
		}
	}

}

func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

// CreateExampleFile creates an example NetCDF file containing only one variable.
func CreateExampleFile(filename string) error {
	// Create a new NetCDF 4 file. The dataset is returned.
	ds, err := netcdf.CreateFile("gopher.nc", netcdf.CLOBBER|netcdf.NETCDF4)
	if err != nil {
		return err
	}
	defer ds.Close()

	// Add the dimensions for our data to the dataset
	dims := make([]netcdf.Dim, 2)
	ht, wd := 5, 4
	dims[0], err = ds.AddDim("height", uint64(ht))
	if err != nil {
		return err
	}
	dims[1], err = ds.AddDim("width", uint64(wd))
	if err != nil {
		return err
	}

	// Add the variable to the dataset that will store our data
	v, err := ds.AddVar("gopher", netcdf.UBYTE, dims)
	if err != nil {
		return err
	}

	// Add a _FillValue to the variable's attributes
	// From C++ netCDF documentation:
	//   With netCDF-4 files, nc_put_att will notice if you are writing a _FillValue attribute,
	//   and will tell the HDF5 layer to use the specified fill value for that variable. With
	//   either classic or netCDF-4 files, a _FillValue attribute will be checked for validity,
	//   to make sure it has only one value and that its type matches the type of the associated
	//   variable.
	if err := v.Attr("_FillValue").WriteUint8s([]uint8{255}); err != nil {
		return err
	}

	// Add an attribute to the variable
	if err := v.Attr("year").WriteInt32s([]int32{2012}); err != nil {
		return err
	}

	// Create the data with the above dimensions and write it to the file.
	gopher := make([]uint8, ht*wd)
	i := 0
	for y := 0; y < ht; y++ {
		for x := 0; x < wd; x++ {
			gopher[i] = uint8(x + y)
			i++
		}
	}
	return v.WriteUint8s(gopher)
}

// ReadExampleFile reads the data in NetCDF file at filename and prints it out.
func ReadExampleFile(filename string) error {
	// Open example file in read-only mode. The dataset is returned.
	ds, err := netcdf.OpenFile(filename, netcdf.NOWRITE)
	if err != nil {
		return err
	}
	defer ds.Close()

	// Get the variable containing our data and read the data from the variable.
	v, err := ds.Var("gopher")
	if err != nil {
		return err
	}

	// Print variable attribute
	year, err := netcdf.GetInt32s(v.Attr("year"))
	if err != nil {
		return err
	}
	fmt.Printf("year = %v\n", year[0])

	// Read data from variable
	gopher, err := netcdf.GetUint8s(v)
	if err != nil {
		return err
	}

	// Get the length of the dimensions of the data.
	dims, err := v.LenDims()
	if err != nil {
		return err
	}

	fmt.Println(dims)

	// Print out the data
	i := 0
	for y := 0; y < int(dims[0]); y++ {
		for x := 0; x < int(dims[1]); x++ {
			fmt.Printf(" %d", gopher[i])
			i++
		}
		fmt.Printf("\n")
	}
	return nil
}

func getVarFloat(ds *netcdf.Dataset, varname string, scaleFactor float64) ([]float64, error) {
	// var d []float64
	d := make([]float64, 1)

	// Get the variable containing our data and read the data from the variable.
	fidv, err := ds.Var(varname)
	if err != nil {
		return d, err
	}

	// Read data from variable
	fid, err := netcdf.GetInt32s(fidv)
	if err != nil {
		return d, err
	}

	dims, err := fidv.LenDims()
	if err != nil {
		return d, err
	}

	for y := 0; y < int(dims[0]); y++ {
		// d[y] = float64(fid[y]) * scaleFactor
		d = append(d, float64(fid[y])*scaleFactor)
	}

	return d, err
}

func getVarInt(ds *netcdf.Dataset, varname string) ([]int32, error) {
	d := make([]int32, 1)

	// Get the variable containing our data and read the data from the variable.
	fidv, err := ds.Var(varname)
	if err != nil {
		return d, err
	}

	// Read data from variable
	fid, err := netcdf.GetInt32s(fidv)
	if err != nil {
		return d, err
	}

	dims, err := fidv.LenDims()
	if err != nil {
		return d, err
	}

	for y := 0; y < int(dims[0]); y++ {
		d = append(d, fid[y])
	}

	return d, err
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
		log.Println("Error building subject IRI for tika triple")
		return buf.String(), err
	}
	pred, err := rdf.NewIRI(p)
	if err != nil {
		log.Println("Error building predicate IRI for tika triple")
		return buf.String(), err
	}
	obj, err := rdf.NewLiteral(o)
	if err != nil {
		log.Println("Error building object literal for tika triple")
		return buf.String(), err
	}

	t := rdf.Triple{Subj: sub, Pred: pred, Obj: obj}
	q := rdf.Quad{t, ctx}

	qs := q.Serialize(rdf.NQuads)
	if s != "" && p != "" && o != "" {
		fmt.Fprintf(buf, "%s", qs)
	}
	return buf.String(), err
}
