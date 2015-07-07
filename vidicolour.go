package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/lazywei/go-opencv/opencv"
)

func saveImages(path int, cap *opencv.Capture, fpm int) {
	img := cap.QueryFrame()
	opencv.SaveImage(fmt.Sprintf("./%d.png", path), img, fpm)
	go colourSample(fmt.Sprintf("./%d.png", path))
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
	fmt.Println(dominantColor(color_index))
}

func dominantColor(colors chan<- [][3]uint8) string {
	return fmt.Sprintf("%v", colors)
}

func convert(r, g, b, _ uint32) (uint8, uint8, uint8) {
	return uint8(r), uint8(g), uint8(b)
}

func main() {

	filename := "mov.mp4"

	cap := opencv.NewFileCapture(filename)
	if cap == nil {
		panic("can not open video")
	}
	defer cap.Release()

	// // Get an image every 5 minutes
	fpm := int(cap.GetProperty(opencv.CV_CAP_PROP_FPS) * (60 * 5))

	var img *opencv.IplImage

	c := make(chan int)

	i := 1
	for {
		cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(fpm*i))
		img = cap.QueryFrame()
		if img != nil {
			go saveImages(i, cap, fpm)
			i++
		} else {
			break
		}
	}

	done := <-c
	fmt.Println(done)
}
