package image

import (
	"testing"
)

func TestBitsInThreesFromBytes(t *testing.T) {
	test(t, []byte{1},      []uint8{0, 0, 2})
	test(t, []byte{129},    []uint8{4, 0, 2})
	test(t, []byte{255},    []uint8{7, 7, 6})
	test(t, []byte{0, 0},   []uint8{0, 0, 0, 0, 0, 0})
	test(t, []byte{1, 129}, []uint8{0, 0, 3, 0, 0, 4})
}

func TestDecode(t *testing.T) {
	var tests = []struct{
		in []byte
		want string
	}{
		{[]byte{}, ""},
		{[]byte{97}, "a"},
		{[]byte{97, 13, 98}, "a\rb"},
	}

	for _, c := range tests {
		got := string(decode(c.in))

		if got != c.want {
			t.Errorf("Wanted %q, got %q\n", c.want, got)
		}
	}
}

func test(t *testing.T, input []byte, expected []uint8) {
	ch := make(chan uint8)

	go bitsInThreesFromBytes(input, ch)

	result := readThreeBits(ch, len(expected))

	if ! same(result, expected) {
		t.Error("%s != %s", result, expected)
	}
}

func readThreeBits(ch chan uint8, howMany int) []uint8 {
	result := make([]uint8, howMany)
	index  := 0

	for value := range ch {
		result[index] = value
		index += 1
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
