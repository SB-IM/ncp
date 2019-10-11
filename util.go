package main

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

