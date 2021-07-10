package main

import(
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"math"
	"strings"
)

var verboseFlag bool = false
const chunkSize = 16

func init() {
	flag.BoolVar(&verboseFlag, "verbose", false, "Print more information")
	flag.Parse()
}

func loadImage(path string) image.Image {
	handle, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer handle.Close()
	img, _, err := image.Decode(handle)
	if err != nil {
		panic(err)
	}
	return img
}

func saveImage(path string, img *image.RGBA) {
	handle, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer handle.Close()
	err = png.Encode(handle, img)
	if err != nil {
		panic(err)
	}
}

func copyChunk(src image.Image, dst *image.RGBA, srcX, srcY, dstX, dstY int) {
	for y := 0; y < chunkSize; y++ {
		for x := 0; x < chunkSize; x++ {
			clr := src.At(srcX + x, srcY + y)
			dst.Set(dstX + x, dstY + y, clr)
		}
	}
}

func process(path string) {
	if verboseFlag {
		fmt.Println("Currently processing", path)
	}
	img := loadImage(path)
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	w, h = w / chunkSize, h / chunkSize

	dim := int(math.Sqrt(float64(w * h)) + 0.5)
	newImg := image.NewRGBA(image.Rect(
		0,
		0,
		int(dim * chunkSize),
		int(dim * chunkSize),
	))

	chunksToCopy := w * h
	srcX, srcY, dstX, dstY := 0, 0, 0, 0
	for i := 0; i < chunksToCopy; i++ {
		if i > 0 && i % w == 0 {
			srcX = 0
			srcY += chunkSize
		}

		if i > 0 && i % dim == 0 {
			dstX = 0
			dstY+= chunkSize
		}

		copyChunk(img, newImg, srcX, srcY, dstX, dstY)

		srcX += chunkSize
		dstX += chunkSize
	}

	saveStr := ""
	index := strings.LastIndex(path, ".")
	if index == -1 {
		saveStr = path + "_square.png"
	} else {
		saveStr = path[:index] + "_square" + path[index:]
	}

	if verboseFlag {
		fmt.Println("Saving results to", saveStr)
	}
	saveImage(saveStr, newImg)
}

func main() {
	targets := flag.Args()
	for _, target := range targets {
		process(target)
	}
}
