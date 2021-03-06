package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	contentDirectoryPattern = `\d\d\d\d-\d\d-\d\d - (.*)`
	todayFormat             = "2006-01-02"
)

var (
	reIsImageFile        = regexp.MustCompile(`(?i)\.(jpg|jpeg|png)$`)
	reIsContentDirectory = regexp.MustCompile(`\d\d\d\d-\d\d-\d\d - (.*)`)
	todayDate            = time.Now().Format(todayFormat)
)

type FilterPredicate func(f os.FileInfo) bool

// Filter predicate - only directories with date in names
func FilterOnlyFoldersWithContent() FilterPredicate {
	return func(f os.FileInfo) bool {
		return f.IsDir() && reIsContentDirectory.MatchString(f.Name())
	}
}

// Filter predicate - only directories with TODAY date in names
func FilterOnlyTodayFolders() FilterPredicate {
	return func(f os.FileInfo) bool {
		return f.IsDir() && strings.Contains(f.Name(), todayDate)
	}
}

// Filter predicate - only JPEG/PNG images (exluding '_resized_' prefix)
func FilterImageFiles() FilterPredicate {
	return func(f os.FileInfo) bool {
		name := f.Name()
		return !f.IsDir() && reIsImageFile.MatchString(name) && !strings.Contains(name, resizedPrefix)
	}
}

// Filter - get files names from FileInfo list
func FilterFilesNames(files []os.FileInfo) []string {
	var filesNames []string
	for _, file := range files {
		filesNames = append(filesNames, file.Name())
	}
	return filesNames
}

// Filter - filter list of FileInfo by predicate
func FilterFilesList(files []os.FileInfo, filterPredicate FilterPredicate) []os.FileInfo {
	var filesList []os.FileInfo
	for _, file := range files {
		if filterPredicate(file) {
			filesList = append(filesList, file)
		}
	}
	return filesList
}

// Read directory content and filter files list by predicate
func getFolderContent(folderPath string, filterPredicate FilterPredicate) []os.FileInfo {
	// Read directory
	files, err := ioutil.ReadDir(folderPath)
	checkError(err)
	// Filter by predicate
	return FilterFilesList(files, filterPredicate)
}

// Read title from folder name and images from 'images' subfolder and generate task
func generateTaskFromFolder(folderPath string) Task {
	// Get title from folder name
	titles := reIsContentDirectory.FindStringSubmatch(folderPath)
	if len(titles) < 2 {
		logger.Panic("Title is not found for folder " + folderPath)
	}
	title := titles[1]

	// Find all image files in image folder
	imagesList := getFolderContent(folderPath, FilterImageFiles())
	var imagesFullNamesList []string
	for _, fileName := range FilterFilesNames(imagesList) {
		imagesFullNamesList = append(imagesFullNamesList, folderPath+"/"+fileName)
	}

	// Generate task for processing folder
	return Task{PageTitle: title, Images: imagesFullNamesList}
}
