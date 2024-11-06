package connect

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

func BenchmarkHexAsciiViewFrom(b *testing.B) {
	data := dataForBenchmark()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = HexAsciiViewFrom(data)
	}
}
