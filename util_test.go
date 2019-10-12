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


func TestUcfirst(t *testing.T) {
	if Ucfirst("abc") != "Abc" {
		t.Errorf("Ucfirst")
	}

	if Ucfirst("Abc") != "Abc" {
		t.Errorf("Ucfirst Duplicate")
	}

	if Ucfirst("abc_de") != "Abc_de" {
		t.Errorf("Ucfirst '_' ")
	}
}

func TestLcfirst(t *testing.T) {
	if Lcfirst("ABC") != "aBC" {
		t.Errorf("Lcfirst")
	}

	if Lcfirst("aBC") != "aBC" {
		t.Errorf("Lcfirst Duplicate")
	}

	if Lcfirst("ABC_DE") != "aBC_DE" {
		t.Errorf("Lcfirst '_' ")
	}
}

