// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package helpers

import (
	"fmt"
	"math"
	"strings"
)

func isPrintableASCII(b byte) bool {
	return 32 <= b && b <= 126
}

func writeLineNumber(sb *strings.Builder, i int) {
	sb.Grow(6)
	sb.WriteByte(byte(i>>12) + '0')
	sb.WriteByte(byte(i>>8&0xF) + '0')
	sb.WriteByte(byte(i>>4&0xF) + '0')
	sb.WriteByte(byte(i&0xF) + '0')
	sb.WriteByte(':')
	sb.WriteByte(' ')
}

func writeSizeOfData(sb *strings.Builder, length int) {
	sb.Grow(20)
	sb.WriteString("Size: ")
	maxLog := math.Log10(float64(length))
	maxPower := int(math.Ceil(maxLog))
	for i := maxPower; i > 0; i-- {
		power := int(math.Pow(10, float64(i)))
		value := (length / power) % 10
		if value != 0 {
			sb.WriteByte(byte(value) + '0')
		}
	}
	sb.WriteByte(byte(length%10) + '0')
	sb.WriteString(" bytes\n")
}

// This is overly optimized function was just used as example for myself
// benchmarking. Was no real point in making it more optimized cause it
// is used just for console data display.
func HexAsciiViewFrom(data []byte) string {
	const bytesPerRow = 16
	var sb strings.Builder
	length := len(data)
	sb.Grow(length*5 + 20)

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

		writeLineNumber(&sb, i)
		sb.Write(hexPart)
		sb.WriteByte(' ')
		sb.Write(asciiPart)
		sb.WriteByte('\n')
	}
	writeSizeOfData(&sb, length)

	return sb.String()
}

func HexViewFrom(data []byte) string {
	return HexViewFromWithLineSplit(data, 16)
}

func HexViewFromWithLineSplit(data []byte, lineLength int) string {
	var sb strings.Builder
	for i := 0; i < len(data); i += lineLength {
		for j := 0; j < lineLength; j++ {
			if i+j < len(data) {
				sb.WriteByte("0123456789abcdef"[data[i+j]>>4])
				sb.WriteByte("0123456789abcdef"[data[i+j]&0xF])
			} else {
				sb.WriteString("  ")
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func HexStringFromInt32(i int32) string {
	bytes := []byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i),
	}
	return HexViewFrom(bytes)
}

func ShowAsHexAndAscii(data []byte) {
	fmt.Println(HexAsciiViewFrom(data))
}

func ShowAsHexView(data []byte) {
	fmt.Println(HexViewFrom(data))
}
