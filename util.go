package main

import (
	"bufio"
	"bytes"
	"unicode"
)

type DuplicateFilter struct {
	Msg string
}

func (last *DuplicateFilter) Put(m string) string {
	if last.Msg == m {
		return ""
	} else {
		last.Msg = m
		return m
	}
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// This function mainly solves the case where the number of bytes in a single line is greater than 4096
func readLine(reader *bufio.Reader) ([]byte, error) {
	var buffer bytes.Buffer
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return buffer.Bytes(), err
		}
		buffer.Write(line)
		if !isPrefix {
			break
		}
	}
	return buffer.Bytes(), nil
}
