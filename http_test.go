package main

import (
	"net/http"
	"strings"
	"testing"
)

var access_id, secret_key = "2", "cJXWEPknyfzAZPOmQX6/mbpGqSuzxYD9aezcDHezy4tVK7U94vfxLEObcC9yD4TRfBVddg/ir7XzDDDTn7GFMA=="
var date, mac = "Mon, 17 Jun 2019 03:10:03 GMT", "SuX8+f6zt5TxYoTe0cuo9PALBbY="

func Test_ValidMAC(t *testing.T) {
	str := "GET,,,/ncp/v1/plans/12/get_map," + date

	if !ValidMAC([]byte(str), []byte(secret_key), mac) {
		t.Errorf("ValidMAC Error")
	}

	if ValidMAC([]byte(str), []byte(secret_key), "") {
		t.Errorf("ValidMAC Error")
	}
}

func Test_signHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ncp/v1/plans/12/get_map", nil)
	req.Header.Add("DATE", date)

	signHeader(access_id, secret_key, req)
	req.Header.Add("Authorization", "APIAuth "+access_id+":"+SignMAC([]byte(getCanonicalString(req)), []byte(secret_key)))

	//t.Log(mac)
	if strings.Split(req.Header.Get("Authorization"), ":")[1] != mac {
		t.Errorf("signHeader Error")
	}
}
