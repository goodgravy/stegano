package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/goodgravy/stegano/image"
	"os"
	"strings"
)

var inputFile string
var outputFile string
var revealing bool

func init() {
	flag.StringVar(&inputFile, "input", "original.png", "Path to the image used as input")
	flag.StringVar(&outputFile, "output", "output.png", "When hiding, we will destructively write the new image to this path")
	flag.BoolVar(&revealing, "reveal", false, "Instead reveal the text hidden in an image")

	flag.Parse()
}

func readTextToHide() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Text to hide: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(text)
}

func main() {
	if revealing {
		text := image.RevealTextInImage(inputFile)
		fmt.Printf("Your text: %q\n", text)
	} else {
		image.HideStringInImage(readTextToHide(), inputFile, outputFile)
	}
}
