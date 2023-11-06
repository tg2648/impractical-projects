package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

// 10 levels of gray
const gscale1 = "@%#*+=-:. "

// 70 levels of gray
const gscale2 = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "

func convertToGrayscale(img image.Image) *image.Gray {
	grayImg := image.NewGray(img.Bounds())
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	return grayImg
}

func writeImgToFile(img image.Image, imgPath string) {
	writer, err := os.Create(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	jpeg.Encode(writer, img, nil)
}

func _getAvgBrightness(img image.Gray, yMin, yMax, xMin, xMax int) uint8 {
	nPixels := (yMax - yMin) * (xMax - xMin)
	total := 0

	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			nPixels += 1
			total += int(img.GrayAt(x, y).Y)
		}
	}

	return uint8(total / nPixels)
}

func getAvgBrightness(img image.Gray, nRow int, nCol int) [][]uint8 {
	avg := make([][]uint8, nRow)
	bounds := img.Bounds()

	perRow := (bounds.Max.Y - bounds.Min.Y) / nRow
	perCol := (bounds.Max.X - bounds.Min.X) / nCol

	for y := 0; y < nRow; y++ {
		avg[y] = make([]uint8, nCol)
		yMin := y * perRow
		yMax := (y + 1) * perRow
		for x := 0; x < nCol; x++ {
			xMin := x * perCol
			xMax := (x + 1) * perCol
			avg[y][x] = _getAvgBrightness(img, yMin, yMax, xMin, xMax)
		}
	}

	return avg
}

func getAscii(avgBrightness [][]uint8, scale string) string {
	var builder strings.Builder

	gradient := (255 / len(scale)) + 1

	for _, row := range avgBrightness {
		for _, val := range row {
			idx := int(val) / gradient
			builder.WriteString(string(scale[idx]))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		log.Fatal("Usage: ascii <image_path> <rows> <columns>")
	}

	imgPath := args[0]
	nRow, err1 := strconv.Atoi(args[1])
	nCol, err2 := strconv.Atoi(args[2])
	if err1 != nil || err2 != nil {
		log.Fatal("Number of columns/rows should be a valid integer")
	}

	reader, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	grayImg := convertToGrayscale(img)
	// imgName := strings.TrimSuffix(path.Base(imgPath), path.Ext(imgPath))
	// writeImgToFile(grayImg, fmt.Sprintf("%v_gray.jpg", imgName))

	avgBrightness := getAvgBrightness(*grayImg, nRow, nCol)
	ascii := getAscii(avgBrightness, gscale2)

	fmt.Println(ascii)
}
