package main

import (
	"code.google.com/p/go-opencv/trunk/opencv"
	"fmt"
	"image/png"
	"log"
	"os"
)

func saveImages(img_path string, img *opencv.IplImage, fpm int) {
	opencv.SaveImage(img_path, img, 0)
	opencv.WaitKey(fpm)
	colourSample(img_path)
}

func colourSample(path string) {
	immg, err := os.Open(path)
	defer immg.Close()

	if err != nil {
		panic(err)
	}

	color_index := make(chan<- [][3]uint8)
	decoded_image, _ := png.Decode(immg)
	max_y := decoded_image.Bounds().Max.Y
	max_x := decoded_image.Bounds().Max.X

	go func() {
		var color_map [][3]uint8
		for i := 0; i <= max_y; i += 5 {
			for y := 0; y <= max_x; y += 50 {
				var r, g, b = convert(decoded_image.At(max_x-y, max_y-i).RGBA())
				var rgb = [3]uint8{r, g, b}
				color_map = append(color_map, rgb)
			}
		}

		color_index <- color_map
	}()

	// Get the dominant color
	dominantColor(color_index)
}

func dominantColor(colors chan<- [][3]uint8) {
	log.Printf("%s", colors)
}

func convert(r, g, b, _ uint32) (uint8, uint8, uint8) {
	return uint8(r), uint8(g), uint8(b)
}

func main() {

	img_path := "sample_images/"

	filename := "movie.mkv"
	if len(os.Args) == 2 {
		filename = os.Args[1]
	} else {
		fmt.Printf("Usage: go run vidicolour.go video_file \n")
		os.Exit(0)
	}

	cap := opencv.NewFileCapture(filename)
	if cap == nil {
		panic("can not open video")
	}
	defer cap.Release()

	// Get an image every 5 minutes
	fpm := int(cap.GetProperty(opencv.CV_CAP_PROP_FPS) * (60 * 5))
	// // First query
	img := cap.QueryFrame()

	i := 1
	for {
		cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(fpm*i))
		img = cap.QueryFrame()
		if img != nil {
			img_path = fmt.Sprintf("sample_images/%d.png", i)
			go saveImages(img_path, img, fpm) // offload to goroutine
			i++
		} else {
			break
		}

	}

}
