// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"encoding/hex"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func textDiff(a, b string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(a, b, true)
	return dmp.DiffPrettyText(diffs)
}

func testData() []byte {
	line0 := "000102030405060708090a0b0c0d0e0f"
	line1 := "101112131415161718191a1b1c1d1e1f"
	line2 := "202122232425262728292a2b2c2d2e" // Intentionally one byte short
	data, err := hex.DecodeString(line0 + line1 + line2)
	if err != nil {
		panic(err)
	}
	return data
}

func TestHexViewEmptyData(t *testing.T) {
	result := HexAsciiViewFrom([]byte{})

	if result != "Size: 0 bytes\n" {
		t.Errorf("Got different output:\n%s", textDiff(result, ""))
	}
}

func TestHexAsciiView(t *testing.T) {
	testData := testData()
	result := HexAsciiViewFrom(testData)

	line1 := "0000: 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f  ................\n"
	line2 := "0010: 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f  ................\n"
	line3 := "0020: 20 21 22 23 24 25 26 27 28 29 2a 2b 2c 2d 2e      !\"#$%&'()*+,-. \n"
	line4 := "Size: 47 bytes\n"
	expected := line1 + line2 + line3 + line4

	if result != expected {
		t.Errorf("Got different output:\n%s", textDiff(result, expected))
	}
}

func TestShowAsHexAndAscii(t *testing.T) {
	ShowAsHexAndAscii(testData())
}
