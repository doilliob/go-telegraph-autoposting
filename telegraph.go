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
	paths, err := client.Upload(images)
	checkError(err)
	return paths
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
