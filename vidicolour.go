package main

import (
	"code.google.com/p/go-opencv/trunk/opencv"
	"fmt"
	"os"
)

func saveImages(img_path string, img *opencv.IplImage, fpm int) {
	opencv.SaveImage(img_path, img, 0)
	opencv.WaitKey(fpm)
}

func main() {

	img_path := "sample_images/"

	filename := nil
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

	fmt.Printf("%s", "Adasdasdsa")

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
