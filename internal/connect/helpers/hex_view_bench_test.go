// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package helpers

import (
	"testing"
)

func dataForBenchmark() []byte {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

func BenchmarkPlusPlus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		value := 2 + 2
		if value != 4 {
			b.Fatal("value is not 4")
		}
	}
}

func BenchmarkHexASCIIViewFrom(b *testing.B) {
	data := dataForBenchmark()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		result := HexASCIIViewFrom(data)
		if result == "" {
			b.Fatal("result is empty")
		}
	}
}
