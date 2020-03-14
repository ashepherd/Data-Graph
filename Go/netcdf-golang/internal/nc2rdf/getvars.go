package nc2rdf

import (
	"github.com/fhs/go-netcdf/netcdf"
)

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
