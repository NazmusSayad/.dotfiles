package utils

import (
	"os"
	"unicode/utf16"
)

func WriteUTF16LE(path string, s string) error {
	r := utf16.Encode([]rune(s))
	b := make([]byte, 2+len(r)*2)
	b[0], b[1] = 0xFF, 0xFE
	for i, v := range r {
		b[2+i*2] = byte(v)
		b[3+i*2] = byte(v >> 8)
	}
	return os.WriteFile(path, b, 0644)
}
