package main

import (
	"github.com/kallydev/telegraph-go"
)

type Task struct {
	PageTitle string
	Images    []string
}

func executeTask(client *telegraph.Client, t Task) {
	logger.Trace("---------------------------------------------------")
	logger.Trace("Prepare task to execution...")
	logger.Trace("Page title: " + t.PageTitle)
	logger.Trace("Images:")
	for _, image := range t.Images {
		logger.Trace(" + " + image)
	}

	logger.Trace(" -- Resizing big images...")
	resizedImages := ResizeFilesListIfNeeded(t.Images)

	logger.Trace("")
	logger.Trace("New images list:")
	for _, image := range resizedImages {
		logger.Trace(" + " + image)
	}

	logger.Trace("")
	logger.Trace("Uploading images...")
	urls := uploadImages(client, resizedImages)
	for _, url := range urls {
		logger.Trace(" -> " + Url(url))
	}

	logger.Trace("")
	logger.Trace("Public page...")
	pageURL := publicPage(client, t.PageTitle, urls)

	logger.Info("")
	logger.Info("---------------------------------")
	logger.Info("TITLE: " + t.PageTitle)
	logger.Info("URL: " + pageURL)
}
