package util

import (
	"testing"
)

func TestDetachTran(t *testing.T) {
	input := `{"weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}}`
	dataMap := DetachTran([]byte(input))

	output_tag := "weather"
	output := `{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}`
	if string(dataMap[output_tag]) != output {
		t.Errorf("input is : " + input)
		t.Errorf(`No match tag : "` + output_tag + `", data : ` + output)
	}
}

func TestDetachTran_multiple(t *testing.T) {
	input := `{"weather":{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780},"battery":{"status":"ok"}}`
	dataMap := DetachTran([]byte(input))

	output_tag := "weather"
	output := `{"WD":0,"WS":0,"T":66115,"RH":426,"Pa":99780}`
	if string(dataMap[output_tag]) != output {
		t.Errorf("input is : " + input)
		t.Errorf(`No match tag : "` + output_tag + `", data : ` + output)
	}
	output2_tag := "battery"
	output2 := `{"status":"ok"}`
	if string(dataMap[output2_tag]) != output2 {
		t.Errorf("input is : " + input)
		t.Errorf(`No match tag : "` + output2_tag + `", data : ` + output2)
	}
}
