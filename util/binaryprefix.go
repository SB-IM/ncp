package util

import (
	"errors"
	"strconv"
)

func BinaryPrefix(raw string) (result int64, err error) {
	if raw == "" {
		err = errors.New("str is empty")
		return
	}
	binaryArr := []byte{'B', 'K', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y'}

	for i, v := range binaryArr {
		if raw[len(raw)-1] == v {
			result, err = strconv.ParseInt(raw[:len(raw)-1], 10, 64)
			for j := 0; j < i; j++ {
				result *= 1024
			}
		}
	}

	// IF not suffix
	if result == 0 {
		result, err = strconv.ParseInt(raw, 10, 64)
	}
	return
}
