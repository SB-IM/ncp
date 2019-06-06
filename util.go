package main

type DuplicateFilter struct {
  msg string
}

func (last *DuplicateFilter) Put(m string) (string) {
  if last.msg == m {
    return ""
  } else {
    last.msg = m
    return m
  }
}

