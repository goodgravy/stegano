package stenoimage

import (
	"unicode/utf8"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func decode(b []byte) []rune {
	var runes []rune
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		runes = append(runes, r)
		i += size
	}
	return runes
}

func bitsInThreesFromBytes(input []byte, threeBitChan chan byte) {
	bitChan := make(chan uint8)

	go bitsFromBytes(input, bitChan)
	chunkBits(3, bitChan, threeBitChan)
}

func chunkBits(chunkSize int, bitChan chan uint8, byteChan chan byte) {
	for {
		current := byte(0)

		for count := chunkSize - 1; count >= 0; count -= 1 {
			bit, more := <-bitChan
			current += bit << uint8(count)
			if ! more {
				byteChan <- current
				close(byteChan)
				return
			}
		}
		byteChan <- current
	}
}

func bitsFromBytes(input []byte, bitChan chan uint8) {
	for _, value := range input {
		for bitIndex := uint8(0); bitIndex < 8; bitIndex += 1 {
			mask := 0xFE >> bitIndex
			bit := value & uint8(mask)
			bitChan <- bit >> (7 - bitIndex)
		}
	}
	close(bitChan)
}

