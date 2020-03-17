package main

import (
	"fmt"
	"log"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

// <pyarrow._parquet.ParquetSchema object at 0x7ffa0a003128>
type ufokn struct {
	ID        string  `parquet:"name=id, type=UTF8, encoding=PLAIN_DICTIONARY"`        // BYTE_ARRAY String
	Type      string  `parquet:"name=type, type=UTF8, encoding=PLAIN_DICTIONARY"`      // BYTE_ARRAY String
	Name      string  `parquet:"name=name, type=UTF8, encoding=PLAIN_DICTIONARY"`      // BYTE_ARRAY String
	Hand      string  `parquet:"name=hand, type=UTF8, encoding=PLAIN_DICTIONARY"`      // DOUBLE
	Offset    string  `parquet:"name=offset, type=UTF8, encoding=PLAIN_DICTIONARY"`    // DOUBLE
	Featureid string  `parquet:"name=featureid, type=UTF8, encoding=PLAIN_DICTIONARY"` // DOUBLE
	X         float64 `parquet:"name=x, type=DOUBLE"`                                  // DOUBLE
	Y         float64 `parquet:"name=y, type=DOUBLE"`                                  // DOUBLE
}

func main() {
	//	var data []*ufokn
	log.Println("Reading file")

	test()

	//all, err := myreadPartialParquet(50835, 0)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//log.Println("Printing all data")
	//for _, a := range all {
	//	fmt.Println(a)
	//}
}
func test() {
	fr, err := local.NewLocalFileReader("jefferson_complete_risk_pts_test.parquet")
	if err != nil {
		log.Println(err)
	}
	log.Println(" -------------------   check mark -1-----------------------")

	pr0, err := reader.NewParquetReader(fr, new(ufokn), 4)
	if err != nil {
		log.Println(err)
	}
	log.Println(" -------------------   check mark 0-----------------------")

	log.Println(pr0.GetNumRows())

	u := make([]*ufokn, pr0.GetNumRows())

	log.Println(" -------------------   check mark 0.5-----------------------")

	if err = pr0.Read(&u); err != nil {
		log.Println(" -------------------   check mark 1-----------------------")
		log.Println(err)
	}

	pr0.ReadStop()

	for i := range u {
		fmt.Print(u[i])
	}

	// // Get a column reader and loop over the rows.
	// pr, err := reader.NewParquetColumnReader(fr, 4) // 4 is the number parallel which I suspect is a go routine count?
	// if err != nil {
	// 	log.Println(err)
	// }

	// num := pr.GetNumRows()
	// d2, _, _, err := pr.ReadColumnByIndex(2, num)
	// d1, _, _, err := pr.ReadColumnByIndex(1, num)
	// if err != nil {
	// 	log.Println(err)
	// }

	// for i := range d2 {
	// 	fmt.Printf("col1: %s   col2: %s \n", d1[i], d2[i])
	// }

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
