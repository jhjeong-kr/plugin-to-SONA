package util

import (
	"encoding/json"
)

func InterfaceToString(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "converting error"
	}
	return string(bytes)
}

func InterfaceToIndenttedString(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return "converting error"
	}
	return string(bytes)
}
