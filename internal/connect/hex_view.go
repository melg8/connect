// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"fmt"
	"strings"
)

func isPrintableASCII(b byte) bool {
	return 32 <= b && b <= 126
}

func HexAsciiViewFrom(data []byte) string {
	const bytesPerRow = 16
	var sb strings.Builder
	length := len(data)
	sb.Grow(length * 5)

	for i := 0; i < length; i += bytesPerRow {
		hexPart := make([]byte, bytesPerRow*3)
		asciiPart := make([]byte, bytesPerRow)

		for j := 0; j < bytesPerRow; j++ {
			pos := j * 3
			if i+j < length {
				b := data[i+j]
				hexPart[pos] = "0123456789abcdef"[b>>4]
				hexPart[pos+1] = "0123456789abcdef"[b&0xF]
				hexPart[pos+2] = ' '
				if isPrintableASCII(b) {
					asciiPart[j] = b
				} else {
					asciiPart[j] = '.'
				}
			} else {
				hexPart[pos] = ' '
				hexPart[pos+1] = ' '
				hexPart[pos+2] = ' '
				asciiPart[j] = ' '
			}
		}

		sb.WriteString(fmt.Sprintf("%04x: ", i))
		sb.Write(hexPart)
		sb.WriteByte(' ')
		sb.Write(asciiPart)
		sb.WriteByte('\n')
	}

	return sb.String()
}

func ShowAsHexAndAscii(data []byte) {
	fmt.Println(HexAsciiViewFrom(data))
}
