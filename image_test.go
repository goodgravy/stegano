package main

import (
	"fmt"
	"testing"
)

func TestBitsFromBytes(t *testing.T) {
	test([]byte{1}, []uint8{0, 0, 0, 0, 0, 0, 0, 1})
	test([]byte{129}, []uint8{1, 0, 0, 0, 0, 0, 0, 1})
	test([]byte{255}, []uint8{1, 1, 1, 1, 1, 1, 1, 1})
	test([]byte{0, 0}, []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	test([]byte{1, 1}, []uint8{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1})
}

func test(input []byte, expected []uint8) (bool, error) {
	ch := make(chan uint8)

	go BitsFromBytes(input, ch)

	result := readBits(ch, len(expected))

	if ! same(result, expected) {
		return false, fmt.Errorf("%s != %s", result, expected)
	}
	return true, nil
}

func readBits(ch chan uint8, howMany int) []uint8 {
	result := make([]uint8, howMany)
	for i := 0; i < len(result); i++ {
		result[i] = <-ch
	}
	return result
}

func same(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}