package utils

import "strings"

func ReplaceLastOcurrence(s, oldChar, newChar string) string {
	index := strings.LastIndex(s, oldChar)
	if index == -1 {
		return s
	}

	return s[:index] + newChar + s[index+len(oldChar):]
}
