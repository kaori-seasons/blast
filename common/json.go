package common

import (
	"bytes"
	"encoding/json"
)

func DecodeJson(body []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(body))
	decoder.UseNumber()

	return decoder.Decode(&v)
}
