package plain

import (
	"bytes"
	"math"
)

// CalculateSimilarity between two strings
func CalculateSimilarity(str1, str2 []byte) float64 {

	if (str1 == nil) || (str2) == nil {
		return 0
	}
	if (len(str1) == 0) || (len(str2) == 0) {
		return 0
	}

	if bytes.Equal(str1, str2) {
		return 1
	}

	stepsToSame := Levenshtein(str1, str2)

	result := (1.0 - (float64(stepsToSame) / float64(maximum(len(str1), len(str2)))))
	return math.Floor(result*100) / 100
}

// Levenshtein calculation
func Levenshtein(str1, str2 []byte) int {

	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]

}

func maximum(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
