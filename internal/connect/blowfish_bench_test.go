package connect

import (
	"testing"
)

func dataForBlowfishBenchmark(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

func BenchmarkBlowFish(b *testing.B) {
	data := dataForBlowfishBenchmark(1000000)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err := DefaultAuthKey().Decrypt(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
