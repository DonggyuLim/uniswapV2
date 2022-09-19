package utils

import (
	"bytes"
	"encoding/gob"
)

// data -> byteSlice
func DataToByte(data interface{}) []byte {
	var result bytes.Buffer
	enc := gob.NewEncoder(&result)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}
