package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/jpeg"
	"os"
	"unicode/utf8"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func BitsInThreesFromBytes(input []byte, threeBitChan chan byte) {
	bitChan := make(chan uint8)

	go bitsFromBytes(input, bitChan)
	chunkBits(3, bitChan, threeBitChan)
}

func BytesFromImage(img *image.RGBA, byteChan chan byte) {
	bitChan := make(chan uint8)

	go bitsFromImage(img, bitChan)
	chunkBits(8, bitChan, byteChan)
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

func bitsFromImage(img *image.RGBA, bitChan chan uint8) {
	for pixOffset := 1; pixOffset < 25; pixOffset += 1 {
		r, g, b, _ := img.At(pixOffset, 1).RGBA()
		bitChan <- uint8(r) & 1
		bitChan <- uint8(g) & 1
		bitChan <- uint8(b) & 1
	}
	close(bitChan)
}

func imageToRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)
	return m
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

func hideBitsInImage(img *image.RGBA, threeBitChan chan uint8) {
	pixCounter := 0
	for threeBit := range threeBitChan {
		rBit := (threeBit & 4) >> 2
		gBit := (threeBit & 2) >> 1
		bBit := threeBit & 1

        y := pixCounter / img.Stride + 1
        x := pixCounter % img.Stride + 1

        r, g, b, _ := img.At(x, y).RGBA()

        img.SetRGBA(x, y, color.RGBA{
        	R: (uint8(r) & 0xFE) | rBit,
        	G: (uint8(g) & 0xFE) | gBit,
        	B: (uint8(b) & 0xFE) | bBit,
        	A: 0xFF,
        })

        pixCounter += 1
	}
}

func decodeFile(filename string) image.Image {
	inFile, err := os.Open(filename)
	check(err)
	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	img, _, err := image.Decode(reader)
	check(err)

	fmt.Println("Read", filename)
	return img
}

func encodePNG(filename string, img image.Image) {
	fo, err := os.Create(filename)
	check(err)
	defer fo.Close()
	defer fo.Sync()

	writer := bufio.NewWriter(fo)
	defer writer.Flush()

	err = png.Encode(writer, img)
	check(err)
	fmt.Println("Wrote to", filename)
}

func main() {
	rgbIm := imageToRGBA(decodeFile("original.jpeg"))

	input := []byte("hello")
	input = append(input, 0)
	threeBitChan := make(chan byte)

	go BitsInThreesFromBytes(input, threeBitChan)
	hideBitsInImage(rgbIm, threeBitChan)

	encodePNG("output.png", rgbIm)

	rgbIm = imageToRGBA(decodeFile("output.png"))

	byteChan := make(chan byte)
	go BytesFromImage(rgbIm, byteChan)

	var bytes []byte

	for v := range byteChan {
		if v == 0 {
			fmt.Println("EOF found")
			break
		}
		bytes = append(bytes, v)
	}

	fmt.Printf("%q\n", string(decode(bytes)))
}
