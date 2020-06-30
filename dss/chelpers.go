package dss

import (
	"C"
	"unsafe"
)

//GoStrings converts an array of type **char in c to a []string in go
//https://stackoverflow.com/questions/36188649/cgo-char-to-slice-string
func GoStrings(n C.int, charList **C.char) []string {
	nStrings := int(n)
	tmpSlice := (*[1 << 30]*C.char)(unsafe.Pointer(charList))[:nStrings:nStrings]
	goStrings := make([]string, nStrings)

	for i, s := range tmpSlice {
		goStrings[i] = C.GoString(s)
	}
	return goStrings
}
