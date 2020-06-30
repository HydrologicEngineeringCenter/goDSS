package dss

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -ljavaheclib -L.
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
	cDSSPaths := C.CString("*")
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
