// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package toauthserver

import "testing"

func BenchmarkRequestGGAuth_ToBytes(b *testing.B) {
	req := &RequestGGAuth{
		SessionID: 1,
		Data1:     2,
		Data2:     3,
		Data3:     4,
		Data4:     5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := req.ToBytes()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRequestGGAuth_ToBytesDirectWriter(b *testing.B) {
	req := &RequestGGAuth{
		SessionID: 1,
		Data1:     2,
		Data2:     3,
		Data3:     4,
		Data4:     5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := req.ToBytesDirectWriter()
		if err != nil {
			b.Fatal(err)
		}
	}
}
