package xunsafe

import (
	"reflect"
	"unsafe"
)

// BytesToString converts a byte slice to a string without memory allocation.
//
// The byte slice must not be used, modified, or reallocated after this call
// since the returned string references the same memory.
func BytesToString(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}

// StringToBytes converts a string to a byte slice without memory allocation.
//
// The string must not be used, modified, or reallocated after this call since the returned byte slice references the same memory.
func StringToBytes(str string) []byte {
	if str == "" {
		return []byte{}
	}

	var slice []byte

	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))

	sliceHeader.Data = stringHeader.Data
	sliceHeader.Len = stringHeader.Len
	sliceHeader.Cap = stringHeader.Len

	return slice
}
