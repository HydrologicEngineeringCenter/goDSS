package dss

import (
	"C"
	"unsafe"
)

//GoStrings converts an array of type **char in c to a []string in go
//https://stackoverflow.com/questions/36188649/cgo-char-to-slice-string
func GoStrings(n C.int, charList **C.char) []string {
	nStrings := int(n)
	// Need to verify 1 << 30 behavior
	tmpSlice := (*[1 << 30]*C.char)(unsafe.Pointer(charList))[:nStrings:nStrings]
	goStrings := make([]string, nStrings)

	for i, s := range tmpSlice {
		goStrings[i] = C.GoString(s)
	}
	return goStrings
}

//GoFloat64s converts an array of type **char in c to a []string in go
func GoFloat64s(n C.int, floatList *C.double) []float64 {
	nFloats := int(n)
	goFloats := make([]float64, nFloats)
	tmpSlice := (*[1 << 30]C.double)(unsafe.Pointer(floatList))[:nFloats:nFloats]
	for i, s := range tmpSlice {
		goFloats[i] = float64(s)
	}
	return goFloats
}
