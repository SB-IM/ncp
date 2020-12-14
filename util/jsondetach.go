package util

import (
	"encoding/json"
)

func DetachTran(raw []byte) map[string][]byte {
	srcMap := make(map[string]*json.RawMessage)
	dstMap := make(map[string][]byte)
	json.Unmarshal(raw, &srcMap)
	for k, v := range srcMap {
		dstMap[k], _ = v.MarshalJSON()
	}
	return dstMap
}
