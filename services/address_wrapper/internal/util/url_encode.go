package util

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-._"
)

func URLEncode(value string) string {
	encoded := ""
	valueCopy := value
	for len(valueCopy) > 0 {
		char, size := utf8.DecodeRuneInString(valueCopy)
		if char == utf8.RuneError {
			break
		}
		if !strings.ContainsRune(alphanumeric, char) {
			var completeHex = fmt.Sprintf("%X", char)
			// left-pad the first byte for 3-byte unicode
			if len(completeHex)%2 == 1 {
				completeHex = "0" + completeHex
			}
			// split the full hex into 2 bytes prefixed with '%'
			var currentBytes = ""
			for _, h := range completeHex {
				currentBytes += string(h)
				if len(currentBytes)%2 == 0 {
					encoded += "%" + currentBytes
					currentBytes = ""
				}
			}
		} else {
			encoded += string(char)
		}
		valueCopy = valueCopy[size:]
	}
	return encoded
}
