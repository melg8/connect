// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package crypt

import (
	"maps"
	"math/rand/v2"
	"slices"
	"sort"
	"testing"
	"time"
)

func UniqRandn1(n int) []int {
	if n <= 0 {
		return []int{}
	}
	unique := make(map[int]struct{}, n)
	value := 0
	for len(unique) < n {
		value = rand.Int()
		unique[value] = struct{}{}
	}
	return slices.Collect(maps.Keys(unique))
}

func UniqRandn2(n int) []int {
	if n <= 0 {
		return nil
	}

	uniq := make([]int, 0, n)
	seen := make(map[int]struct{}, n)
	for len(uniq) < n {
		val := rand.Int()
		if _, exists := seen[val]; exists {
			continue
		}
		seen[val] = struct{}{}
		uniq = append(uniq, val)
	}
	return uniq
}

func UniqRandn3(n int) []int {
	if n <= 0 {
		return nil
	}

	uniq := make([]int, n)
	seen := make(map[int]struct{}, n)
	pos := 0
	var val int
	seenValue := struct{}{}
	for pos < n {
		val = rand.Int()
		if _, exists := seen[val]; exists {
			continue
		}
		seen[val] = seenValue
		uniq[pos] = val
		pos++
	}
	return uniq
}

func Unique[T comparable](xs []T) []T {
	n := len(xs)
	if n == 0 {
		return xs
	}
	j := 0
	for i := 1; i < n; i++ {
		if xs[j] != xs[i] {
			j++
			if j < i {
				xs[j] = xs[i]
				for k := i + 1; k < n; k++ {
					if xs[j] != xs[k] {
						j++
						xs[j] = xs[k]
					}
				}
				break
			}
		}
	}
	xs = xs[0 : j+1]
	return xs
}

func UniqRandn4(n int) []int {
	if n <= 0 {
		return nil
	}

	uniq := make([]int, n)
	generatedN := 0

	for generatedN < n {
		for i := generatedN; i < n; i++ {
			uniq[i] = rand.Int()
		}
		sort.Ints(uniq)
		uniq = Unique(uniq)
		generatedN = len(uniq)
	}

	return uniq
}

type RandomSequenceOfUnique struct {
	mIndex              uint32
	mIntermediateOffset uint32
}

func permuteQPR(x uint32) uint32 {
	const prime = 4294967291
	if x >= prime {
		return x // The 5 integers out of range are mapped to themselves.
	}
	residue := uint32((uint64(x) * uint64(x)) % prime)
	if x <= prime/2 {
		return residue
	}
	return prime - residue
}

func NewRandomSequenceOfUnique(seedBase, seedOffset uint32) RandomSequenceOfUnique {
	return RandomSequenceOfUnique{
		mIndex:              permuteQPR(permuteQPR(seedBase) + 0x682f0161),
		mIntermediateOffset: permuteQPR(permuteQPR(seedOffset) + 0x46790905),
	}
}

func (r *RandomSequenceOfUnique) Next() uint32 {
	result := permuteQPR((permuteQPR(r.mIndex) + r.mIntermediateOffset) ^ 0x5bf03635)
	r.mIndex++
	return result
}

func UniqRandn5(n int) []int {
	if n <= 0 {
		return nil
	}

	uniq := make([]int, n)
	unixTimeSeed := uint32(time.Now().Unix())
	rsu := NewRandomSequenceOfUnique(unixTimeSeed, unixTimeSeed+1)

	for i := 0; i < n; i++ {
		uniq[i] = int(rsu.Next())
	}
	return uniq
}

func benchmarkUniqRandn(b *testing.B, fn func(int) []int) {
	for n := 0; n < b.N; n++ {
		res := fn(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 100_000")
		}
	}
}

func BenchmarkUniqRandn1(b *testing.B) {
	benchmarkUniqRandn(b, UniqRandn1)
}

func BenchmarkUniqRandn2(b *testing.B) {
	benchmarkUniqRandn(b, UniqRandn2)
}

func BenchmarkUniqRandn3(b *testing.B) {
	benchmarkUniqRandn(b, UniqRandn3)
}

func BenchmarkUniqRandn4(b *testing.B) {
	benchmarkUniqRandn(b, UniqRandn4)
}

func BenchmarkUniqRandn5(b *testing.B) {
	benchmarkUniqRandn(b, UniqRandn5)
}

func randn(n int) []int {
	if n <= 0 {
		return nil
	}

	uniq := make([]int, n)
	for i := 0; i < n; i++ {
		uniq[i] = rand.Int()
	}
	return uniq
}

func BenchmarkRandn(b *testing.B) {
	benchmarkUniqRandn(b, randn)
}
