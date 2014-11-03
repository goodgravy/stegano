package stenoimage

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/jpeg"
	"os"
)

func bytesFromImage(img *image.RGBA, byteChan chan byte) {
	bitChan := make(chan uint8)

	go bitsFromImage(img, bitChan)
	chunkBits(8, bitChan, byteChan)
}

func bitsFromImage(img *image.RGBA, bitChan chan uint8) {
	for pixOffset := 1; ; pixOffset += 1 {
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

func hideBitsInImage(img *image.RGBA, threeBitChan chan uint8) {
	pixCounter := 0
	for threeBit := range threeBitChan {
		rBit := (threeBit & 4) >> 2
		gBit := (threeBit & 2) >> 1
		bBit := threeBit & 1

        y := pixCounter / img.Stride + 1
        x := pixCounter % img.Stride + 1

        r, g, b, a := img.At(x, y).RGBA()

        img.SetRGBA(x, y, color.RGBA{
        	R: (uint8(r) & 0xFE) | rBit,
        	G: (uint8(g) & 0xFE) | gBit,
        	B: (uint8(b) & 0xFE) | bBit,
        	A: uint8(a),
        })

        pixCounter += 1
	}
	bounds := img.Bounds()
	min := bounds.Min
	max := bounds.Max
	fmt.Printf("Used %d / %d pixels to hide message\n", pixCounter, (max.X - min.X) * (max.Y - min.Y))
}

func decodeImage(filename string) image.Image {
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

func RevealTextInImage(inputFile string) string {
	rgbIm := imageToRGBA(decodeImage(inputFile))

	byteChan := make(chan byte)
	go bytesFromImage(rgbIm, byteChan)

	var bytes []byte

	for v := range byteChan {
		if v == 0 {
			break
		}
		bytes = append(bytes, v)
	}

	return string(decode(bytes))
}

func HideStringInImage(toHide, inputFile, outputFile string) {
	rgbIm := imageToRGBA(decodeImage(inputFile))

	input := []byte(toHide)
	input = append(input, 0)
	threeBitChan := make(chan byte)

	go bitsInThreesFromBytes(input, threeBitChan)
	hideBitsInImage(rgbIm, threeBitChan)

	encodePNG(outputFile, rgbIm)
}