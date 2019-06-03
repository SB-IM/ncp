package main

import (
  "testing"
)

func TestDuplicateFilter(t *testing.T) {
  m := DuplicateFilter{}
  m.Put("aa")

  if m.Put("aa") != "" {
    t.Errorf("No Duplicate Filter")
  }

  if m.Put("iaa") != "iaa" {
    t.Errorf("Excessive Filter")
  }

}

