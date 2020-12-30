package main

import (
	"github.com/kallydev/telegraph-go"
)

func getTelegraphConnection(acountName string) (*telegraph.Client, *telegraph.Account) {
	client, err := telegraph.NewClient("", nil)
	checkError(err)

	account, err := client.CreateAccount(acountName, nil)
	checkError(err)

	client.AccessToken = account.AccessToken
	return client, account
}

func uploadImages(client *telegraph.Client, images []string) []string {
	var urls []string
	for _, imageName := range images {
		logger.Trace("Uploading: " + imageName)
		imageURL, err := client.Upload([]string{imageName})
		checkError(err)
		logger.Trace("Image URL: " + imageURL[0])
		urls = append(urls, imageURL[0])
	}
	return urls
}

func publicPage(client *telegraph.Client, title string, imagesURLs []string) string {
	// Convert image URLs to NodeElement graph
	var imagesNodes []telegraph.Node
	for _, url := range imagesURLs {
		newElement := telegraph.NodeElement{
			Tag:      "img",
			Attrs:    map[string]string{"src": Url(url)},
			Children: []telegraph.Node{},
		}
		imagesNodes = append(imagesNodes, newElement)
	}
	// Public page
	content := []telegraph.Node{
		title,
		telegraph.NodeElement{
			Tag:      "p",
			Attrs:    map[string]string{},
			Children: imagesNodes,
		},
	}
	page, err := client.CreatePage(title, content, nil)
	checkError(err)

	return page.URL
}
