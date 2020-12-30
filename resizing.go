package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"regexp"

	"github.com/nfnt/resize"
)

const (
	resizedPrefix = "_resized_"
	resizeStep    = 0.05
)

var (
	reIsFileJPEG = regexp.MustCompile(`(?i)\.(jpeg|jpg)$`)
	reIsFilePNG  = regexp.MustCompile(`(?i)\.png$`)
	reExtension  = regexp.MustCompile(`(?i)\.[a-z]+$`)
)

func resizeFile(fileName string, img image.Image, newWidth uint) {
	m := resize.Resize(newWidth, 0, img, resize.Lanczos3)
	out, err := os.Create(fileName)
	checkError(err)
	defer out.Close()
	jpeg.Encode(out, m, nil)
}

func recursiveResizeImage(img image.Image, width int, fileName string, currentPercent float64, limitSizeBytes int64) {
	newPercent := currentPercent - resizeStep
	if newPercent <= 0 {
		panic(fmt.Sprintf("Image %s can not be resized correctly", fileName))
	}
	newWidth := uint(float64(width) * newPercent)
	logger.Trace(fmt.Sprintf(" -- Resizing image to new width %d", newWidth))
	resizeFile(fileName, img, newWidth)

	// Check if size is not exceed limit
	fileinfo, err := os.Stat(fileName)
	checkError(err)
	if fileinfo.Size() <= limitSizeBytes {
		logger.Trace(" -- Resizing is not needed.")
		return
	}
	logger.Trace(fmt.Sprintf(" -- File size %d is not enough", fileinfo.Size()))
	recursiveResizeImage(img, width, fileName, newPercent, limitSizeBytes)
}

type Decoder func(r io.Reader) (image.Image, error)

func getDecoder(filePath string) Decoder {
	if !reIsFileJPEG.MatchString(filePath) && !reIsFilePNG.MatchString(filePath) {
		logger.Panic(fmt.Sprintf("File %s is not JPEG/PNG file!", filePath))
	}
	if reIsFileJPEG.MatchString(filePath) {
		return jpeg.Decode
	}
	return png.Decode
}

func loadImageFromFile(filePath string) image.Image {
	// Open file
	file, err := os.Open(filePath)
	checkError(err)
	defer file.Close()

	// Decode image from file
	img, err := (getDecoder(filePath))(file)
	checkError(err)
	return img
}

func resizeImageIfNeeded(path string) string {
	// If file smaller - return his name
	maxSizeBytes := 1024 * 1024 * int64(configuration.MaxImageSizeMb)
	stat, err := os.Stat(path)
	checkError(err)
	if stat.Size() <= maxSizeBytes {
		return path
	}

	// Resize image while new image size
	img := loadImageFromFile(path)
	resizedFileName := reExtension.ReplaceAllString(path, "") + resizedPrefix + ".jpg"
	width := img.Bounds().Max.X - img.Bounds().Min.X
	percent := 1.0
	// Logging
	logger.Trace(" -- ")
	logger.Trace(fmt.Sprintf(" [!!!] File %s size %d bytes is more than maximum limit %d", path, stat.Size(), maxSizeBytes))
	logger.Trace(" -- Resized file name: " + resizedFileName)
	logger.Trace(fmt.Sprintf(" -- Base file width: %d", width))
	// Resize image
	recursiveResizeImage(img, width, resizedFileName, percent, maxSizeBytes)

	// Return resized filename
	return resizedFileName
}

// Resize each file if needed
func ResizeFilesListIfNeeded(files []string) []string {
	var filesList []string
	for _, file := range files {
		filesList = append(filesList, resizeImageIfNeeded(file))
	}
	return filesList
}
