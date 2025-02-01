package main

import (
	"strings"
	"unicode"
)

func sanitizeFileName(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ':':
			return '-'
		case '/', '\\', '<', '>', '"', '*', '|', '?':
			return -1
		default:
			if !unicode.IsPrint(r) {
				return -1
			}

			return r
		}
	}, s)
}
