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

//GoFloats converts an array of type **char in c to a []string in go
//https://stackoverflow.com/questions/36188649/cgo-char-to-slice-string
func GoFloats(n C.int, floatList *C.double) []float32 {
	nFloats := int(n)
	goFloats := make([]float32, nFloats)
	tmpSlice := (*[1 << 28]C.double)(unsafe.Pointer(floatList))[:nFloats:nFloats]
	for i, s := range tmpSlice {
		goFloats[i] = float32(s)
	}
	return goFloats
}
