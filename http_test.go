package main

import (
  "testing"
)

func Test_ValidMAC(t *testing.T) {
  //access_id, secret_key := "2", "cJXWEPknyfzAZPOmQX6/mbpGqSuzxYD9aezcDHezy4tVK7U94vfxLEObcC9yD4TRfBVddg/ir7XzDDDTn7GFMA=="
  secret_key := "cJXWEPknyfzAZPOmQX6/mbpGqSuzxYD9aezcDHezy4tVK7U94vfxLEObcC9yD4TRfBVddg/ir7XzDDDTn7GFMA=="

  date := "Mon, 17 Jun 2019 03:10:03 GMT"
  str := "GET,,,/ncp/v1/plans/12/get_map," + date

  if !ValidMAC([]byte(str), []byte(secret_key), "SuX8+f6zt5TxYoTe0cuo9PALBbY=") {
    t.Errorf("ValidMAC Error")
  }

  if ValidMAC([]byte(str), []byte(secret_key), "") {
    t.Errorf("ValidMAC Error")
  }
}

