package main

import (
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

