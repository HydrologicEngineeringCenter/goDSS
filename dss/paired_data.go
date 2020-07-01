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

//HelloWorld tests clean opening and closing of a dss file
func HelloWorld(filePath string) {
	ifltab := C.longlong(250)
	cPath := C.CString(filePath)

	fmt.Println("Hello DSS!\n")
	C.zopen(&ifltab, cPath)
	defer C.zclose(&ifltab)
}

//ReadCatalogue tests listing paths in dss file
func ReadCatalogue(filePath string) {
	ifltab := C.longlong(250)
	cPath := C.CString(filePath)
	cDSSPaths := C.CString("/*/*/*/*/*/*/")
	cField := C.int(1)

	dssFilename := C.mallocAndCopy(cPath)

	C.zopen(&ifltab, dssFilename)
	defer C.zclose(&ifltab)

	catStruct := C.zstructCatalogNew()
	defer C.zstructFree(unsafe.Pointer(catStruct))

	nPaths := C.zcatalog(&ifltab, cDSSPaths, catStruct, cField)
	// fmt.Println(nPaths == catStruct.numberPathnames)

	pathNameList := GoStrings(nPaths, catStruct.pathnameList)

	fmt.Println("File List: \n")
	for i := 0; i < int(nPaths); i++ {
		fmt.Println(i, pathNameList[i])
	}
	fmt.Println("\nClosing File.")

}

//ReadTimeSeries tests listing paths in dss file
func ReadTimeSeries(filePath string) {
	ifltab := C.longlong(250)
	cPath := C.CString(filePath)
	cDSSPaths := C.CString("/*/*/*FLOW*/*/*/*/")
	cField := C.int(1)
	cDate := C.CString("startDate")
	cDateLength := C.int(13)
	cTime := C.CString("startTime")
	cTimeLength := C.int(10)

	C.zopen(&ifltab, cPath)
	defer C.zclose(&ifltab)

	catStruct := C.zstructCatalogNew()
	defer C.zstructFree(unsafe.Pointer(catStruct))

	nPaths := C.zcatalog(&ifltab, cDSSPaths, catStruct, cField)

	pathNameList := GoStrings(nPaths, catStruct.pathnameList)

	// Add selection criteria here
	pathIndex := C.int(1)
	dssRecord := pathNameList[pathIndex]

	tss2 := C.zstructTsNew(C.CString(dssRecord))
	defer C.zstructFree(unsafe.Pointer(tss2))

	C.ztsRetrieve(&ifltab, tss2, -1, 2, 0)

	valueTime := tss2.startTimeSeconds / tss2.timeGranularitySeconds

	nValues := tss2.numberValues
	doubleValues := GoFloats(nValues, tss2.doubleValues)

	for i := 0; i < int(nValues); i++ {
		C.getDateAndTime(C.int(valueTime),
			tss2.timeGranularitySeconds,
			tss2.startJulianDate,
			cDate,
			cDateLength,
			cTime,
			cTimeLength)

		valueTime += tss2.timeIntervalSeconds / tss2.timeGranularitySeconds
		fmt.Println(i, C.GoString(cDate), C.GoString(cTime), doubleValues[i])
	}

	fmt.Println("\nClosing File.")

}
