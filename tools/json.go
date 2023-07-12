package tools

import "encoding/json"

func MarshalJSON(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func UnmarshalJSON[T any](jsonStr string) (T, error) {
	var tar T
	err := json.Unmarshal([]byte(jsonStr), &tar)
	return tar, err
}
