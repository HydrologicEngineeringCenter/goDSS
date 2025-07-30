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
	"errors"
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
type RegularTimeSeries struct {
	Times  []time.Time //Times  []int //- how do i convert this dataset
	Values []float64
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
	ts, err := dssFile.ReadRegularTimeSeries("//livingston_s030/FLOW/*/1Hour/RUN:SST/")
	if err != nil {
		panic(err)
	}
	ts.Print()
}
func InitDssFile(filepath string) (DssFile, error) {
	cFilepath := C.CString(filepath)
	//defer C.free(unsafe.Pointer(cFilepath))

	var cDss *C.dss_file
	file := DssFile{}

	ret := C.hec_dss_open(cFilepath, &cDss)

	if ret != 0 {
		return file, fmt.Errorf("error opening DSS file, check path and file validity")
	}

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
	if int(a) == 0 {
		fmt.Println("catalog is empty")
		return DssCatalog{}, errors.New("empty dss catalog")
	}
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

func (d DssFile) ReadRegularTimeSeries(pathname string) (RegularTimeSeries, error) {
	cPathname := C.CString(pathname)
	cFullSet := C.int(1) //expose as option eventually?
	cStartDate := C.int(0)
	cStartTime := C.int(0)
	cEndDate := C.int(0)
	cEndTime := C.int(0)

	C.hec_dss_tsGetDateTimeRange(d.fileHandle, cPathname, cFullSet, &cStartDate, &cStartTime, &cEndDate, &cEndTime)

	startTimeGo := julianToGoTime(cStartDate, cStartTime)
	startDateString := strings.ToUpper(startTimeGo.Format("02Jan2006"))
	startTimeString := startTimeGo.Format("15:04:05")
	fmt.Println(startDateString + " " + startTimeString)
	endTimeGo := julianToGoTime(cEndDate, cEndTime)
	endDateString := strings.ToUpper(endTimeGo.Format("02Jan2006"))
	endTimeString := endTimeGo.Format("15:04:05")
	fmt.Println(endDateString + " " + endTimeString)
	bufferLength := 40
	units := make([]byte, bufferLength)
	cUnits := (*C.char)(unsafe.Pointer(&units[0]))
	mytype := make([]byte, bufferLength)
	cType := (*C.char)(unsafe.Pointer(&mytype[0]))
	cLength := C.int(bufferLength)

	timeLength := 10

	cStartDateChar := C.CString(startDateString)
	cEndDateChar := C.CString(endDateString)
	cStartTimeChar := C.CString(startTimeString)
	cEndTimeChar := C.CString(endTimeString)

	cNumPeriods := C.int(0)
	cQualitySize := C.int(0)
	getSizesErr := C.hec_dss_tsGetSizes(d.fileHandle, cPathname, cStartDateChar, cStartTimeChar, cEndDateChar, cEndTimeChar, &cNumPeriods, &cQualitySize)
	if int(getSizesErr) != 0 {
		return RegularTimeSeries{}, fmt.Errorf("could not determine dimensions of pathname %v", pathname)
	}

	times := make([]int, int(cNumPeriods))
	cTimes := (*C.int)(unsafe.Pointer(&times[0]))
	vals := make([]float64, int(cNumPeriods))
	cVals := (*C.double)(unsafe.Pointer(&vals[0]))
	cJulian := C.int(0)
	qualities := make([]int, int(cNumPeriods))
	cQualities := (*C.int)(unsafe.Pointer(&qualities[0]))
	cArraySize := cNumPeriods
	cGranularity := C.int(3600)
	timezone := make([]byte, timeLength)
	cTimezone := (*C.char)(unsafe.Pointer(&timezone[0]))
	cValsRead := C.int(0)
	response := C.hec_dss_tsRetrieve(d.fileHandle, cPathname, cStartDateChar, cStartTimeChar, cEndDateChar, cEndTimeChar, cTimes, cVals, cArraySize, &cValsRead, cQualities, cQualitySize, &cJulian, &cGranularity, cUnits, cLength, cType, cLength, cTimezone, cLength)
	timesGo := intTimeArrayToGoTimeArray(times, 3600, startTimeGo)
	vals = vals[0:len(timesGo)]
	if int(response) != 0 {
		return RegularTimeSeries{}, fmt.Errorf("could not read regular timeseries at pathname %v", pathname)
	}
	return RegularTimeSeries{timesGo, vals}, nil
}
func secondsToHours(seconds int) int {
	return seconds / 60 / 60
}
func julianToGoTime(cDate C.int, cTime C.int) time.Time {
	cYear := C.int(0)
	cMonth := C.int(0)
	cDay := C.int(0)
	C.hec_dss_julianToYearMonthDay(cDate, &cYear, &cMonth, &cDay)
	Month := time.Month(int(cMonth))
	Hour := secondsToHours(int(cTime))
	t := time.Date(int(cYear), Month, int(cDay), Hour, 0, 0, 0, time.UTC)
	if Hour == 24 {
		t = t.Add(time.Nanosecond * -1)
	}
	return t
}
func (catalog *DssCatalog) RemoveDatesFromCatalog() {
	for i, p := range catalog.PathNames {
		parts := strings.Split(p, "/")
		d := parts[4]
		newp := strings.Replace(p, d, "*", -1)
		catalog.PathNames[i] = newp
	}
}

func intTimeArrayToGoTimeArray(times []int, granularity int, startDateTime time.Time) []time.Time {
	timesGo := make([]time.Time, 0)

	var delta time.Duration
	//d := time.Duration(t)
	switch granularity {
	case 1:
		delta = time.Second // * d
	case 60:
		delta = time.Minute // * d
	case 3600:
		delta = time.Hour //* d
	case 86400:
		delta = time.Hour * 24 //* d
	}
	currentTime := startDateTime
	for i := range times {
		if i == 0 {
			timesGo = append(timesGo, startDateTime)
		} else {

			currentTime = currentTime.Add(delta)
			timesGo = append(timesGo, currentTime)
		}
	}
	return timesGo
}
func (ts RegularTimeSeries) Print() {
	fmt.Printf("Times,Values\n")
	for i, v := range ts.Values {
		fmt.Printf("%v,%v\n", ts.Times[i].Format("02Jan2006 15:04:05"), v)
	}
}
