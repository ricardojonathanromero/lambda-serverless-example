package utils

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
)

// ToString converts interface to string/*
func ToString(i interface{}) string {
	m, _ := json.Marshal(i)
	return string(m)
}

// StringToInt converts string to int64/*
func StringToInt(value string) int64 {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return int64(v)
}

func EncodeStr(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
