package main

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -ljavaheclib -L.
// #include "dss/headers/heclib7.h"
// #include "dss/headers/hecdss7.h"
import "C"

import (
	"fmt"
	"os"

	"github.com/HydrologicEngineeringCenter/goDSS/dss"
)

/*
NOTE: Make sure the LD_LIBRARY_PATH is set prior to compiling, and the following
files are in the dss directory:
	1. heclib.a
	2. libjavaHeclib.so

Example:
		cd to  GOPATH/github.com/HydrologicEngineeringCenter/goDSS/dss
		export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:=${PWD}
*/

func main() {

	var usageWarning string = "./hello_dss dss/data/G14.dss"
	filePath := os.Args[1]

	if len(os.Args) != 2 {
		fmt.Println(usageWarning)
	}

	dssContents := dss.ReadCatalogue(filePath)

	// Print all paths and all time series from the test file
	for i := 0; i < len(dssContents); i++ {
		recordPath := dssContents[i]
		tSeries := dss.ReadTimeSeries(filePath, recordPath)
		fmt.Println(len(tSeries))

		// Seeing an issue on print with results here
		// Not seeing this when printing from inside function
		for j := 0; j < len(tSeries); j++ {
			fmt.Println(recordPath, tSeries[j])
		}

	}

}
