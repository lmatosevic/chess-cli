package utils

import (
	"encoding/json"
	"io"
)

func ParseJson[T any](content io.Reader) (T, error) {
	var model T
	decoder := json.NewDecoder(content)
	err := decoder.Decode(&model)
	return model, err
}

func ConvertJson(model any) (string, error) {
	content, err := json.Marshal(model)
	return string(content), err
}
