package main

/*
#cgo LDFLAGS: -lhecdss -L.
#include <dlfcn.h>
#include <stdlib.h>
#include <stdio.h>
#include </dss/lib/hec-dss-7-IU-15/heclib/hecdss/hecdss.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

type DssFile struct {
	fileHandle *C.dss_file
}
type DssCatalog struct {
	PathNames   []string
	RecordTypes []int
}

func main() {

	myfile, err := InitDssFile("/workspaces/goDss/SST.dss")

	if err != nil {
		panic(err)
	}
	defer myfile.Close()
	catalog, err := myfile.ReadCatalog()
	if err != nil {
		panic(err)
	}
	fmt.Println(catalog.PathNames[0])
	/*err = myfile.ReadTimeSeries(catalog.PathNames[0])
	if err != nil {
		panic(err)
	}*/
}
func InitDssFile(filepath string) (DssFile, error) {
	cFilepath := C.CString(filepath)
	//defer C.free(unsafe.Pointer(cFilepath))

	var cDss *C.dss_file
	file := DssFile{}

	ret := C.hec_dss_open(cFilepath, &cDss)

	fmt.Printf("hec_dss_open returned %d\n", int(ret))
	if ret != 0 {
		return file, fmt.Errorf("error opening DSS file, check path and file validity")
	}
	fmt.Printf("success! DSS handle = %p\n", cDss)
	ver := C.hec_dss_getVersion(cDss)
	fmt.Printf("success! DSS version = %v\n", int(ver))

	file.fileHandle = cDss
	return file, nil
}
func (d DssFile) Close() {
	C.hec_dss_close(d.fileHandle)
}

func (d DssFile) ReadCatalog() (DssCatalog, error) {
	cRecordCount := C.hec_dss_record_count(d.fileHandle)
	pathNames := make([]byte, cRecordCount*394)
	cPathBuffer := (*C.char)(unsafe.Pointer(&pathNames[0]))
	cFilter := C.CString("/*/*/*/*/*/*/") //pathpartswithwildcards
	recordTypes := make([]int, cRecordCount)
	cRecordTypes := (*C.int)(unsafe.Pointer(&recordTypes[0]))

	cPathBufferItemSize := C.int(394) //394 defined by  hecdss.c hec_dss_CONSTANT_MAX_PATH_SIZE
	a := C.hec_dss_catalog(d.fileHandle, cPathBuffer, cRecordTypes, cFilter, cRecordCount, cPathBufferItemSize)
	fmt.Print(a)
	stringPathNames := []string{}
	for i := 0; i < int(cRecordCount); i++ {
		bdata := pathNames[i*394 : (i+1)*394]
		sdata := string(bdata)
		sdata = strings.TrimRight(sdata, "\x00")
		stringPathNames = append(stringPathNames, sdata)
	}

	return DssCatalog{
		PathNames:   stringPathNames,
		RecordTypes: recordTypes,
	}, nil
}

/*
func (d DssFile) ReadTimeSeries(pathname string) error {
	cPathname := C.CString(pathname)
	unitLength := 10
	units := make([]byte, unitLength)
	cUnits := (*C.char)(unsafe.Pointer(&units[0]))
	mytype := make([]byte, unitLength)
	cType := (*C.char)(unsafe.Pointer(&mytype[0]))
	cLength := C.int(unitLength)
	C.hec_dss_tsRetrieveInfo(d.fileHandle, cPathname, cUnits, cLength, cType, cLength)
		//timeLength := 10
		//startDate := make([]byte, timeLength)
		//cStartDate := (*C.char)(unsafe.Pointer(&startDate[0]))
		//startTime := make([]byte, timeLength)
		//cStartTime := (*C.char)(unsafe.Pointer(&startTime[0]))
		//cType := (*C.char)(unsafe.Pointer(&mytype[0]))
		//cLength := C.int(unitLength)

	cStartDate := C.CString("")
	cStartTime := C.CString("")
	cEndDate := C.CString("")
	cEndTime := C.CString("")
	numvals := C.int(0)
	qualitySize := C.int(0)

	C.hec_dss_tsGetSizes(d.fileHandle, cPathname, cStartDate, cStartTime, cEndDate, cEndTime, &numvals, &qualitySize)

	return nil
}

*/
