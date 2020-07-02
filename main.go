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

// NOTE: Make sure the LD_LIBRARY_PATH is set prior to compiling, and the following
// files are in the dss directory:
// 	1. heclib.a
// 	2. libjavaHeclib.so

// Example:
// 		cd to  GOPATH/github.com/HydrologicEngineeringCenter/goDSS/dss
// 		export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:=${PWD}
//
// Quick Snips:
// go build main.go
// go build main.go && ./main "dss/data/G14.dss"
// rm main && go build main.go && ./main "dss/data/G14.dss"
func main() {

	var usageWarning string = "./hello_dss dss/data/G14.dss"
	filePath := os.Args[1]

	if len(os.Args) != 2 {
		fmt.Println(usageWarning)
	}

	dssContents := dss.ReadCatalogue(filePath)

	// Print all paths and all time series from the test file to json
	for i := 0; i < len(dssContents); i++ {
		recordPath := dssContents[i]

		jsonFileName := fmt.Sprintf("%d_In_Function.json", i)
		tSeries := make([]dss.TimeSeries, 0)
		tSeries = dss.ReadTimeSeries(filePath, recordPath, jsonFileName)
		fmt.Println(tSeries)

		// jsonFileName2 := fmt.Sprintf("%d_FunctionReturn.json", i)
		// jsonOutput, _ := json.Marshal(tSeries)
		// _ = ioutil.WriteFile(jsonFileName2, tSeries, 0644)

	}

}
