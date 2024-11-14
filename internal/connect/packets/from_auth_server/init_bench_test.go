package from_auth_server

import (
	"testing"
)

func BenchmarkPlusPlus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		value := 2 + 2
		if value != 4 {
			b.Fatal("value is not 4")
		}
	}
}

func dataForBenchmark() []byte {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

func BenchmarkInitPacketParsing(b *testing.B) {
	data := dataForBenchmark()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		packet, err := NewInitPacketFromBytes(data)
		if err != nil {
			panic(err)
		}
		if packet.sessionId != 50462976 {
			b.Fatal("wrong session id")
		}
	}
}
