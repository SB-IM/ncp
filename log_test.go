package main

import (
	"testing"
)

func Test_logGroupNew(t *testing.T) {
	configLog := &ConfigLog{
		Env:   "development",
		Path:  "log/",
		Level: "debug",
		Type: map[string]string{
			"test": "TEST",
		},
	}

	logGroup, err := logGroupNew(configLog)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(logGroup.logs) != 1 {
		t.Errorf("Sum Error")
	}

	if logGroup.Get("test") == nil {
		t.Errorf("No Key 'test'")
	}
}
