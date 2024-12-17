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

func uniqRandn1(n int) []int {
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

func BenchmarkUniqRandn(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := uniqRandn1(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 100_000")
		}
	}
}

func uniqRandn2(n int) []int {
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

func BenchmarkUniqRandn2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := uniqRandn2(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 100_000")
		}
	}
}

func uniqRandn3(n int) []int {
	if n <= 0 {
		return nil
	}

	uniq := make([]int, n)
	seen := make(map[int]struct{}, n)
	pos := 0
	for pos < n {
		val := rand.Int()
		if _, exists := seen[val]; exists {
			continue
		}
		seen[val] = struct{}{}
		uniq[pos] = val
		pos++
	}
	return uniq
}

func BenchmarkUniqRandn3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := uniqRandn3(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 10_000")
		}
	}
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
	for n := 0; n < b.N; n++ {
		res := randn(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 100_000")
		}
	}
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

func uniqRandn5(n int) []int {
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

func BenchmarkUniqRandn5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := uniqRandn5(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 100_000")
		}
	}
}

// func partition(arr []int, low, high int) ([]int, int) {
// 	pivot := arr[high]
// 	i := low
// 	for j := low; j < high; j++ {
// 		if arr[j] < pivot {
// 			arr[i], arr[j] = arr[j], arr[i]
// 			i++
// 		}
// 	}
// 	arr[i], arr[high] = arr[high], arr[i]
// 	return arr, i
// }

// func quickSort(arr []int, low, high int) []int {
// 	if low < high {
// 		var p int
// 		arr, p = partition(arr, low, high)
// 		arr = quickSort(arr, low, p-1)
// 		arr = quickSort(arr, p+1, high)
// 	}
// 	return arr
// }

// func quickSortStart(arr []int) []int {
// 	return quickSort(arr, 0, len(arr)-1)
// }

// func quickSort(arr []int, low int, high int) {
// 	if low < high {
// 		pi := partion(arr, low, high)

// 		// Recursively sort elements before partition and after partition
// 		quickSort(arr, low, pi-1)
// 		quickSort(arr, pi+1, high)
// 	}
// }

// func partion(arr []int, low int, high int) int {
// 	pivot := arr[high]
// 	i := low - 1

// 	for j := low; j < high; j++ {
// 		if arr[j] < pivot {
// 			i++

// 			arr[i], arr[j] = arr[j], arr[i]
// 		}
// 	}

// 	arr[i+1], arr[high] = arr[high], arr[i+1]
// 	return i + 1
// }

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

func uniqRandn4(n int) []int {
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

func BenchmarkUniqRandn4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := uniqRandn4(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 100_000")
		}
	}
}

func uniqRandn10(n int) []int {
	return rand.Perm(n)
}

func BenchmarkUniqRandn10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := uniqRandn10(100_000)
		if len(res) != 100_000 {
			b.Fatal("length is not 100_000")
		}
	}
}
