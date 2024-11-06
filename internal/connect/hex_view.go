// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import "fmt"

func HexAsciiViewFrom(data []byte) string {
	bytesInRow := 16
	result := ""

	for i := 0; i < len(data); i += bytesInRow {
		// Append the offset
		result += fmt.Sprintf("%04x: ", i)

		// Append hex values
		for j := 0; j < bytesInRow; j++ {
			if (i + j) < len(data) {
				result += fmt.Sprintf("%02x ", data[i+j])
			} else {
				result += "   "
			}
		}

		// Append ASCII representation
		result += " "
		for j := 0; j < bytesInRow; j++ {
			if (i + j) < len(data) {
				// Printable ASCII range
				if data[i+j] >= 32 && data[i+j] <= 126 {
					result += string(rune(data[i+j]))
				} else {
					result += "."
				}
			}
		}

		result += "\n"
	}

	return result
}

func ShowAsHexAndAscii(data []byte) {
	fmt.Println(HexAsciiViewFrom(data))
}

func HexViewFrom(data []byte) string {
	bytesInRow := 16
	result := ""

	for i := 0; i < len(data); i += bytesInRow {
		// Append the offset
		result += fmt.Sprintf("%04x: ", i)

		// Append hex values
		for j := 0; j < bytesInRow; j++ {
			if (i + j) < len(data) {
				result += fmt.Sprintf("%02x ", data[i+j])
			} else {
				result += "   "
			}
		}

		result += "\n"
	}

	return result
}

func ShowAsHex(data []byte) {
	fmt.Println(HexViewFrom(data))
}
