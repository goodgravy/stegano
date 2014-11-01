package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/color"
	"os"
)

type StenoImage struct{
	original image.Image
	bitChannel chan uint8
}

type StenoDecoder struct{
	modified image.Image
	bitChannel chan uint8	
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (steno StenoImage) ColorModel() color.Model {
	return steno.original.ColorModel()
}

func (steno StenoImage) Bounds() image.Rectangle {
	return steno.original.Bounds()
}

func (steno StenoImage) At(x, y int) color.Color {
	originalColor := steno.original.At(x, y)

	bit, more := <-steno.bitChannel

	if ! more {
		return originalColor
	}

	red, green, blue, alpha := trimmedColor(originalColor)

	return color.RGBA{
		R: red | bit,
		G: green,
		B: blue,
		A: alpha,
	}
}

func trimmedColor(originalColor color.Color) (uint8, uint8, uint8, uint8) {
	fullRed, fullGreen, fullBlue, fullAlpha := originalColor.RGBA()
	mask := uint8(0xFE)

	return uint8(fullRed) & mask,
		uint8(fullGreen) & mask,
		uint8(fullBlue) & mask,
		uint8(fullAlpha)
}

func BitsFromBytes(bytes []byte, ch chan uint8) {
	for _, value := range bytes {
		for bitIndex := uint8(0); bitIndex < 8; bitIndex += 1 {
			mask := 128 >> bitIndex
			bit := value & uint8(mask)
			ch <- bit >> (7 - bitIndex)
		}
	}
	close(ch)
}

func BitsFromImage(img image.Image, bitChannel chan uint8) {
	// TODO
}

func main() {
    fi, err := os.Open("/tmp/original.jpeg")
    if err != nil {
        panic(err)
    }
    defer func() {
        err := fi.Close()
        check(err)
    }()

    fo, err := os.Create("output.jpeg")
    check(err)
    defer func() {
        err := fo.Close()
        check(err)
    }()

	reader := bufio.NewReader(fi)
    writer := bufio.NewWriter(fo)

	img, _, err := image.Decode(reader)
	check(err)

	bitChannel := make(chan uint8)
	go BitsFromBytes([]byte("hello"), bitChannel)

	steno := StenoImage{
		original: img,
		bitChannel: bitChannel,
	}

    err = jpeg.Encode(writer, steno, nil)
    check(err)

    fmt.Println("wrote encoded file")














    fi, err = os.Open("output.jpeg")
    if err != nil {
        panic(err)
    }
    defer func() {
        err := fi.Close()
        check(err)
    }()

	reader = bufio.NewReader(fi)

	img, _, err = image.Decode(reader)
	check(err)

	bitChannel = make(chan uint8)
	go BitsFromImage(img, bitChannel)

	for v := range bitChannel {
		fmt.Println(v)
	}

	fmt.Println("read encoded file")
}