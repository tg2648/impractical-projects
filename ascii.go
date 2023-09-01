package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("Provide image path")
	}

	img_path := args[0]
	img_name := strings.TrimSuffix(path.Base(img_path), path.Ext(img_path))

	reader, err := os.Open(img_path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	gray_img := image.NewGray(img.Bounds())
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// r, g, b, a := img.At(x, y).RGBA();
			gray_img.Set(x, y, img.At(x, y))
		}
	}

	// img.ColorModel().Convert(color.Gray)
	// fmt.Println(bounds)
	// fmt.Println(img.ColorModel())

	writer, err := os.Create(fmt.Sprintf("%v_gray.jpg", img_name))
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	jpeg.Encode(writer, gray_img, nil)
}
