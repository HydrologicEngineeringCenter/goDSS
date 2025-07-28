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
	"time"
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
	//fmt.Println(catalog.PathNames[0])
	catalog.RemoveDatesFromCatalog()
	//fmt.Println(catalog.PathNames[0])
	err = dssFile.ReadTimeSeries("//livingston_s030/FLOW/*/1Hour/RUN:SST/")
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


		C.hec_dss_tsRetrieveInfo(d.fileHandle, cPathname, cUnits, cLength, cType, cLength)
			//timeLength := 10
			//startDate := make([]byte, timeLength)
			//cStartDate := (*C.char)(unsafe.Pointer(&startDate[0]))
			//startTime := make([]byte, timeLength)
			//cStartTime := (*C.char)(unsafe.Pointer(&startTime[0]))
			//cType := (*C.char)(unsafe.Pointer(&mytype[0]))
			//cLength := C.int(bufferLength)
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
	startMonth := time.Month(int(cStartMonth))
	startTimeGo := time.Date(int(cStartYear), startMonth, int(cStartDay), secondsToHours(int(cStartTime)), 0, 0, 0, time.UTC)
	startDateString := strings.ToUpper(startTimeGo.Format("02Jan2006"))
	startTimeString := startTimeGo.Format("15:04:05")
	fmt.Println(startDateString)
	endMonth := time.Month(int(cEndMonth))
	endTimeGo := time.Date(int(cEndYear), endMonth, int(cEndDay), secondsToHours(int(cEndTime)), 0, 0, 0, time.UTC)
	if endTimeGo.Day() != int(cEndDay) {
		//subtract 1 minute?
		endTimeGo = endTimeGo.Add(time.Minute * -1)
	}
	endDateString := strings.ToUpper(endTimeGo.Format("02Jan2006"))
	endTimeString := endTimeGo.Format("15:04:05")
	fmt.Println(endDateString + " " + endTimeString)
	fmt.Printf("Starting at %v year, %v month, %v day, %v hours\n", int(cStartYear), int(cStartMonth), int(cStartDay), secondsToHours(int(cStartTime)))
	fmt.Printf("Ending at %v year, %v month, %v day, %v hours\n", int(cEndYear), int(cEndMonth), int(cEndDay), secondsToHours(int(cEndTime)))

	bufferLength := 40
	units := make([]byte, bufferLength)
	cUnits := (*C.char)(unsafe.Pointer(&units[0]))
	mytype := make([]byte, bufferLength)
	cType := (*C.char)(unsafe.Pointer(&mytype[0]))
	cLength := C.int(bufferLength)

	timeLength := 10
	//startDate := make([]byte, timeLength)
	cStartDateChar := C.CString(startDateString) //(*C.char)(unsafe.Pointer(&startDateString[0]))
	cEndDateChar := C.CString(endDateString)
	cStartTimeChar := C.CString(startTimeString)
	cEndTimeChar := C.CString(endTimeString)

	cNumPeriods2 := C.int(0)
	cQualitySize := C.int(0)
	getSizesErr := C.hec_dss_tsGetSizes(d.fileHandle, cPathname, cStartDateChar, cStartTimeChar, cEndDateChar, cEndTimeChar, &cNumPeriods2, &cQualitySize)
	fmt.Println(getSizesErr)
	fmt.Println(int(cNumPeriods2))
	fmt.Println(int(cQualitySize))
	//startTime := make([]byte, timeLength)
	//cStartTimeChar := (*C.char)(unsafe.Pointer(&startTime[0]))
	//endDate := make([]byte, timeLength)
	//cEndDateChar := (*C.char)(unsafe.Pointer(&endDateString[0]))
	//endTime := make([]byte, timeLength)
	//cEndTimeChar := (*C.char)(unsafe.Pointer(&endTime[0]))
	times := make([]int, int(cNumPeriods2))
	cTimes := (*C.int)(unsafe.Pointer(&times[0]))
	vals := make([]float64, int(cNumPeriods2))
	cVals := (*C.double)(unsafe.Pointer(&vals[0]))
	cJulian := C.int(0)
	qualities := make([]int, int(cNumPeriods2))
	cQualities := (*C.int)(unsafe.Pointer(&qualities[0]))
	cArraySize := cNumPeriods2
	cGranularity := C.int(3600)
	timezone := make([]byte, timeLength)
	cTimezone := (*C.char)(unsafe.Pointer(&timezone[0]))
	cValsRead := C.int(0)
	response := C.hec_dss_tsRetrieve(d.fileHandle, cPathname, cStartDateChar, cStartTimeChar, cEndDateChar, cEndTimeChar, cTimes, cVals, cArraySize, &cValsRead, cQualities, cQualitySize, &cJulian, &cGranularity, cUnits, cLength, cType, cLength, cTimezone, cLength)
	fmt.Println(vals)
	fmt.Println(response)
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
