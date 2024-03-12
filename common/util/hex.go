package util

import "strings"

func IsHex(value string) bool {
	if strings.Contains(value, "0x") || strings.Contains(value, "0X") {
		return true
	}
	return false
}

func Rm0x(value string) string {
	value = strings.Replace(value, "0x", "", 1)
	value = strings.Replace(value, "0X", "", 1)
	return value
}
