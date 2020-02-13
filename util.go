package main

import (
	"bufio"
	"unicode"
)

type DuplicateFilter struct {
  Msg string
}

func (last *DuplicateFilter) Put(m string) (string) {
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
// Reference: https://www.jianshu.com/p/5b686a9a8eed
func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}

