package xunsafe

import (
	"unsafe"
)

// BytesToString converts a byte slice to a string without memory allocation.
//
// The byte slice must not be used, modified, or reallocated after this call
// since the returned string references the same memory.
func BytesToString(slice []byte) string {
	return unsafe.String(unsafe.SliceData(slice), len(slice))
}

// StringToBytes converts a string to a byte slice without memory allocation.
//
// The string must not be used, modified, or reallocated after this call since the returned byte slice references the same memory.
func StringToBytes(str string) []byte {
	if str == "" {
		return []byte{}
	}

	return unsafe.Slice(unsafe.StringData(str), len(str))
}
