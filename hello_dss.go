package main

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -ljavaheclib -L.
// #include "headers/heclib7.h"
// #include "headers/hecdss7.h"
import "C"

import (
	"fmt"
	"os"
)

// NOTE: Make sure the LD_LIBRARY_PATH is set prior to compiling.
// Example:
// 		export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/mnt/c/Users/slawler/go/src/github.com/HydrologicEngineeringCenter/goDSS
func main() {
	warningMessage := fmt.Sprintf(
		`No DSS file path found. Example usage:\n
            ./hello_dss data/G14.dss`)

	filePath := os.Args[1]

	if len(os.Args) != 1 {
		fmt.Println(warningMessage)
	}

	ifltab := C.longlong(250)
	cPath := C.CString(filePath)

	fmt.Println("Hello DSS!\n")
	C.zopen(&ifltab, cPath)
	C.zclose(&ifltab)
}
