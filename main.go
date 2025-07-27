package godss

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

func demo() {

	dssFile, err := InitDssFile("/workspaces/goDss/SST.dss")

	if err != nil {
		panic(err)
	}
	defer dssFile.Close()
	catalog, err := dssFile.ReadCatalog("/*/*/*/*/*/*/")
	if err != nil {
		panic(err)
	}
	fmt.Println(catalog.PathNames[0])
	catalog.RemoveDatesFromCatalog()
	fmt.Println(catalog.PathNames[0])
	err = dssFile.ReadTimeSeries(catalog.PathNames[0])
	if err != nil {
		panic(err)
	}
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

func (d DssFile) ReadCatalog(filter string) (DssCatalog, error) {
	cRecordCount := C.hec_dss_record_count(d.fileHandle)
	pathNames := make([]byte, cRecordCount*394)
	cPathBuffer := (*C.char)(unsafe.Pointer(&pathNames[0]))
	cFilter := C.CString(filter) //pathpartswithwildcards
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

func (d DssFile) ReadTimeSeries(pathname string) error {
	cPathname := C.CString(pathname)

	/*
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
	*/
	cFullSet := C.int(1)
	cStartDate := C.int(0)
	cStartTime := C.int(0)
	cEndDate := C.int(0)
	cEndTime := C.int(0)
	//numvals := C.int(0)
	//qualitySize := C.int(0)

	C.hec_dss_tsGetDateTimeRange(d.fileHandle, cPathname, cFullSet, &cStartDate, &cStartTime, &cEndDate, &cEndTime)
	fmt.Println(int(cStartDate))
	fmt.Println(int(cStartTime))
	fmt.Println(int(cEndDate))
	fmt.Println(int(cEndTime))

	cInteval := C.int(3600)
	numPeriods := C.hec_dss_numberPeriods(cInteval, cStartDate, cStartTime, cEndDate, cEndTime)
	fmt.Println(numPeriods)

	cStartYear := C.int(0)
	cStartMonth := C.int(0)
	cStartDay := C.int(0)
	cEndYear := C.int(0)
	cEndMonth := C.int(0)
	cEndDay := C.int(0)
	C.hec_dss_julianToYearMonthDay(cStartDate, &cStartYear, &cStartMonth, &cStartDay)
	C.hec_dss_julianToYearMonthDay(cEndDate, &cEndYear, &cEndMonth, &cEndDay)

	fmt.Printf("Starting at %v year, %v month, %v day, %v hours\n", int(cStartYear), int(cStartMonth), int(cStartDay), secondsToHours(int(cStartTime)))
	fmt.Printf("Ending at %v year, %v month, %v day, %v hours\n", int(cEndYear), int(cEndMonth), int(cEndDay), secondsToHours(int(cEndTime)))

	return nil
}
func secondsToHours(seconds int) int {
	return seconds / 60 / 60
}
func (catalog *DssCatalog) RemoveDatesFromCatalog() {
	for i, p := range catalog.PathNames {
		parts := strings.Split(p, "/")
		d := parts[4]
		newp := strings.Replace(p, d, "*", -1)
		catalog.PathNames[i] = newp
	}
}
