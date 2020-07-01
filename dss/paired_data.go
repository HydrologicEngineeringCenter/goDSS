package dss

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -ljavaheclib -L.
// #include <stdio.h>
// #include "headers/heclib7.h"
// #include "headers/hecdss7.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// TimeSeries is a simple container for timeseries data
type TimeSeries struct {
	date   string
	time   string
	value  float64
	status int
}

//HelloWorld tests clean opening and closing of a dss file
func HelloWorld(filePath string) {
	ifltab := C.longlong(250)
	cPath := C.CString(filePath)

	fmt.Println("Hello DSS!\n")
	C.zopen(&ifltab, cPath)
	defer C.zclose(&ifltab)
}

//ReadCatalogue tests listing paths in dss file
func ReadCatalogue(filePath string) []string {
	ifltab := C.longlong(250)
	cPath := C.CString(filePath)

	// Add functionality for selection criteria here
	cDSSPaths := C.CString("/*/*/*/*/*/*/")
	cField := C.int(1)

	C.zopen(&ifltab, cPath)
	defer C.zclose(&ifltab)

	catStruct := C.zstructCatalogNew()
	defer C.zstructFree(unsafe.Pointer(catStruct))

	nPaths := C.zcatalog(&ifltab, cDSSPaths, catStruct, cField)

	return GoStrings(nPaths, catStruct.pathnameList)
}

//ReadTimeSeries tests listing paths in dss file
func ReadTimeSeries(filePath string, record string) []TimeSeries {
	ifltab := C.longlong(250)
	cFilePath := C.CString(filePath)
	recordPath := C.CString(record)

	// Move most or all of these to constants
	cDate := C.CString("startDate")
	cDateLength := C.int(13) // 13 Characters for a date string
	cTime := C.CString("startTime")
	cTimeLength := C.int(10) // 10 Characters for a time string

	C.zopen(&ifltab, cFilePath)
	defer C.zclose(&ifltab)

	catStruct := C.zstructCatalogNew()
	defer C.zstructFree(unsafe.Pointer(catStruct))

	tSeries := C.zstructTsNew(recordPath)
	defer C.zstructFree(unsafe.Pointer(tSeries))

	// https://www.hec.usace.army.mil/confluence/dsscprogrammer/dss-progammers-guide-for-c/time-series-functions
	C.ztsRetrieve(&ifltab, tSeries, -1, 2, 0)

	nValues := tSeries.numberValues
	doubleValues := GoFloat64s(nValues, tSeries.doubleValues)

	// Create holding container for output
	var tsData []TimeSeries

	valueTime := tSeries.startTimeSeconds / tSeries.timeGranularitySeconds

	// This needs to come out of this function
	for i := 0; i < int(nValues); i++ {

		status := C.getDateAndTime(C.int(valueTime),
			tSeries.timeGranularitySeconds,
			tSeries.startJulianDate,
			cDate, cDateLength, cTime, cTimeLength)

		valueTime += tSeries.timeIntervalSeconds / tSeries.timeGranularitySeconds

		ts := TimeSeries{
			date:   C.GoString(cDate),
			time:   C.GoString(cTime),
			value:  doubleValues[i],
			status: int(status)}

		tsData = append(tsData, ts)
	}

	fmt.Println("\nClosing File.")
	return tsData
}
