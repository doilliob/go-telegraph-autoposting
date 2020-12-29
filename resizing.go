package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
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
	logger.Trace(fmt.Sprintf(" -- Resizing image to new width %d", newWidth))
	m := resize.Resize(newWidth, 0, img, resize.Lanczos3)
	out, err := os.Create(fileName)
	checkError(err)
	defer out.Close()
	jpeg.Encode(out, m, nil)
}

func resizeImageIfNeeded(path string) string {
	// If file smaller - return his name
	maxSizeBytes := 1024 * 1024 * int64(configuration.MaxImageSizeMb)
	stat, err := os.Stat(path)
	checkError(err)
	if stat.Size() <= maxSizeBytes {
		return path
	}

	// Else resize image to needed size
	logger.Trace(" -- ")
	logger.Trace(fmt.Sprintf(" [!!!] File %s size %d bytes is more than maximum limit %d", path, stat.Size(), maxSizeBytes))

	// Read image from file
	file, err := os.Open(path)
	checkError(err)
	var img image.Image
	if reIsFileJPEG.MatchString(path) {
		img, err = jpeg.Decode(file)
	} else {
		if reIsFilePNG.MatchString(path) {
			img, err = png.Decode(file)
		} else {
			logger.Panic(fmt.Sprintf("File %s is not JPEG/PNG file!", path))
		}
	}
	checkError(err)
	file.Close()

	// Resize image while new image size
	resizedFileName := reExtension.ReplaceAllString(path, "") + resizedPrefix + ".jpg"
	logger.Trace(" -- Resized file name: " + resizedFileName)
	width := img.Bounds().Max.X - img.Bounds().Min.X
	logger.Trace(fmt.Sprintf(" -- Base file width: %d", width))
	percent := 0.99
	newWidth := uint(float64(width) * percent)
	resizeFile(resizedFileName, img, newWidth)

	for fileinfo, err := os.Stat(resizedFileName); err != nil && fileinfo.Size() > maxSizeBytes; {
		logger.Trace(fmt.Sprintf(" -- File size %d is not enough", fileinfo.Size()))
		percent = percent - resizeStep
		newWidth = uint(float64(width) * percent)
		resizeFile(resizedFileName, img, newWidth)
	}
	logger.Trace(" -- Resizing complete.")

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
